package obdii

import (
	"log"
	"strconv"
	"time"

	"github.com/tarm/serial"
)

const (
	debugSerialPortName               = "OBD-II_DEBUG"
	serialPortName                    = "/dev/ttyUSB0"
	serialPortBaudRate                = 9_600
	serialPortReadTimeoutMilliseconds = 250
	serialPortReadBufferSizeBytes     = 8
)

// Represents an OBD-II serial connection
type Instance struct {
	port         *serial.Port
	readBuffer   []byte
	ecuRequests  chan EcuRequest
	ecuResponses chan EcuResponse
}

// Connect establishes a connection with a serial port and opens communication channels
//
// The application will exit with code 1 if the connection could not be established
func (instance *Instance) Connect(isDebug bool) {
	config := serial.Config{
		Name:        serialPortName,
		ReadTimeout: time.Millisecond * serialPortReadTimeoutMilliseconds,
		Baud:        serialPortBaudRate,
	}

	// Debug mode uses a different serial port name to make it possible to run an OBD-II emulator
	// on a local development machine: https://github.com/Ircama/ELM327-emulator
	//
	// If you use a Windows machine, you should use a null-modem emulator such as:
	// https://sourceforge.net/projects/com0com/
	if isDebug {
		config.Name = debugSerialPortName
	}

	log.Println("Opening serial port connection to port \"" + config.Name + "\"")
	port, error := serial.OpenPort(&config)

	if error != nil {
		log.Fatalln("Unable to open serial port connection: " + error.Error())
	}

	log.Println("Successfully connected to serial port")
	instance.port = port
	instance.readBuffer = make([]byte, serialPortReadBufferSizeBytes)
	instance.ecuRequests = make(chan EcuRequest)
	instance.ecuResponses = make(chan EcuResponse)

	go instance.process()
}

// Disconnect closes the active serial port and closes communication channels
func (instance *Instance) Disconnect() {
	instance.port.Close()
	close(instance.ecuRequests)
	close(instance.ecuResponses)
}

// Get a write-only channel to communicate with the ECU
func (instance *Instance) GetEcuRequestsChannel() chan<- EcuRequest {
	return instance.ecuRequests
}

// Get a read-only channel to listen to the ECU
func (instance *Instance) GetEcuResponsesChannel() <-chan EcuResponse {
	return instance.ecuResponses
}

// OBD-II message processing loop
func (instance *Instance) process() {
	for ecuRequest := range instance.ecuRequests {
		log.Println("Received ECU request", ecuRequest)

		response := EcuResponse{}

		// TODO: add support for sending multiple PIDs in a single command
		for _, pid := range ecuRequest.Pids {
			request := pid.Mode + pid.Pid
			data, err := instance.SendAndWaitForAnswer([]byte(request), pid.ResponseSizeInBytes)

			if err != nil {
				log.Println("Failed to retrieve response from ECU:", err)
				continue
			}

			// TODO: Find a way to parse PIDs
			response.Pids = append(response.Pids, PidResponse{
				Mode: pid.Mode,
				Pid:  pid.Pid,
				Data: string(data),
			})
		}

		instance.ecuResponses <- response
	}
}

// Send an OBD-II command to the ECU and wait for a response
//
// This operation will block until data is received, or until the operation times out
func (instance *Instance) SendAndWaitForAnswer(command []byte, expectedBytes int) ([]byte, error) {
	_, err := instance.port.Write(command)

	// Unable to write to port
	if err != nil {
		log.Println("Failed to write to serial port:", err)
		return nil, err
	}

	data := []byte{}

	for {
		size, _ := instance.port.Read(instance.readBuffer)
		data = append(data, instance.readBuffer[:size]...)

		if size == 0 || (len(data) >= expectedBytes) {
			break
		}
	}

	log.Println("Read " + strconv.Itoa(len(data)) + " bytes from serial port")
	return data, nil
}
