package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tarm/serial"
)

const SerialPortName = "/dev/ttyUSB0"
const SerialPortBaudRate = 9_600

const ContentTypeHeader = "Content-Type"
const MimeTypeText = "text/plain"
const MimeTypeJson = "application/json"

const ObdiiBufferSizeBytes = 64
const ObdiiReadTimeoutMilliseconds = 1_500
const ObdiiReadTimeoutDefaultResponse = "NO DATA AVAILABLE"

type empty = struct{}

type systemStatusResponse struct {
	CpuTemperature float32
	Timestamp      uint64
}

type sendCommandRequest struct {
	Command string
}

type commandResponse struct {
	Data string
}

func handleSystemStatusRequest(response http.ResponseWriter, request *http.Request) {
	data := systemStatusResponse{
		CpuTemperature: -1.0,
		Timestamp:      uint64(time.Now().UnixMilli()),
	}

	json, error := json.Marshal(data)

	if error != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Header().Add(ContentTypeHeader, MimeTypeText)

		response.Write([]byte(error.Error()))
	} else {
		response.Header().Add(ContentTypeHeader, MimeTypeJson)
		response.Write(json)
	}
}

func handleShutdownRequest(response http.ResponseWriter, request *http.Request, quit chan empty) {
	// Only allow DELETE requests
	if request.Method != http.MethodDelete {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	quit <- empty{}
}

func handleObdiiCommandRequest(response http.ResponseWriter, request *http.Request, sender, receiver chan string) {
	// Only allow POST requests
	if request.Method != http.MethodPost {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, error := io.ReadAll(request.Body)

	// No body available
	if error != nil {
		log.Println(error)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	// Body available but unparsable
	var commandRequestBody sendCommandRequest
	if error := json.Unmarshal(bodyBytes, &commandRequestBody); error != nil {
		log.Println(error)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	sender <- commandRequestBody.Command

	for {
		select {
		case result := <-receiver:
			data := commandResponse{Data: result}
			json, error := json.Marshal(data)

			if error != nil {
				log.Println(error)
				response.WriteHeader(http.StatusInternalServerError)
				return
			}

			response.Write(json)
			return
		default:
			time.Sleep(ObdiiReadTimeoutMilliseconds)
		}
	}
}

func poll(sender, receiver chan string, quit chan empty) {
	defer close(sender)
	defer close(quit)
	defer log.Println("Polling loop closed")

	log.Println("Opening serial port...")
	config := serial.Config{
		Name:        SerialPortName,
		ReadTimeout: time.Millisecond * ObdiiReadTimeoutMilliseconds,
		Baud:        SerialPortBaudRate}

	port, error := serial.OpenPort(&config)

	if error != nil {
		log.Fatalln("Cannot open serial port: " + error.Error())
		return
	}

	defer port.Close()

	readBuffer := make([]byte, ObdiiBufferSizeBytes)

	for {
		select {
		case command := <-sender:
			log.Printf("Writing command: %q", command)

			if writeCommandToPort(command, port) {
				log.Println("Waiting for response from port...")
				receiver <- readResponseFromPort(readBuffer, port)
			}
		case <-quit:
			log.Println("Quit signal for polling loop received")
			return
		}
	}
}

func writeCommandToPort(command string, port *serial.Port) bool {
	_, error := port.Write([]byte(command))

	if error != nil {
		log.Println(error)
		return false
	}

	return true
}

func readResponseFromPort(buffer []byte, port *serial.Port) string {
	size, _ := port.Read(buffer)

	if size > 0 {
		log.Print("Port returned raw byte data: ")
		log.Println(buffer[:size])

		var responseBuilder strings.Builder

		// Convert each byte into a hexadecimal character
		for i := 0; i < size; i++ {
			hex := strconv.FormatUint(uint64(buffer[i]), 16)
			responseBuilder.WriteString(hex)
		}

		log.Println("Raw data as formatted response: " + responseBuilder.String())
		return responseBuilder.String()
	}

	log.Println("Port did not respond with any data")
	return ObdiiReadTimeoutDefaultResponse
}

// Reference: https://en.wikipedia.org/wiki/OBD-II_PIDs#Service_01_PID_00
func getSupportedPidsFromEcuResponse(response string, startPid uint8) (map[uint8]bool, bool) {
	supportedPids := make(map[uint8]bool)

	value, error := strconv.ParseUint(response, 16, 32)

	if error != nil {
		log.Println("Unable to parse encoded PIDs into value: " + error.Error())
		return supportedPids, true
	}

	encodedPids := uint32(value)

	// Decrementing loop used here to account for big-endianness of the ECU's response
	for index := 31; index >= 0; index-- {
		isBitSet := encodedPids&(1<<index) != 0
		supportedPids[uint8(index)+startPid] = isBitSet
	}

	return supportedPids, false
}

func main() {
	sender := make(chan string)
	receiver := make(chan string)
	quit := make(chan empty)

	log.Println("Starting polling loop")
	go poll(sender, receiver, quit)

	http.HandleFunc("/api/v1/system/status", handleSystemStatusRequest)
	http.HandleFunc("/api/v1/system", func(response http.ResponseWriter, request *http.Request) {
		handleShutdownRequest(response, request, quit)
	})

	http.HandleFunc("/api/v1/obdii/command", func(response http.ResponseWriter, request *http.Request) {
		handleObdiiCommandRequest(response, request, sender, receiver)
	})

	log.Println("Starting HTTP server on port 80")
	log.Fatal(http.ListenAndServe(":80", nil))
}
