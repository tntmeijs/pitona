package obdii

import (
	"log"
	"strconv"
	"time"

	"github.com/tarm/serial"
	"github.com/tntmeijs/pitona/server/utility"
)

const (
	serialPortName                    = "/dev/ttyUSB0"
	serialPortBaudRate                = 9_600
	serialPortReadTimeoutMilliseconds = 5_000
	serialPortReadBufferSizeBytes     = 4_096
)

// Represents an OBD-II serial connection
type Instance struct {
	port       *serial.Port
	readBuffer []byte
}

// Connect establishes a connection with a serial port
//
// The application will exit with code 1 if the connection could not be established
func (instance *Instance) Connect() {
	config := serial.Config{
		Name:        serialPortName,
		ReadTimeout: time.Millisecond * serialPortReadTimeoutMilliseconds,
		Baud:        serialPortBaudRate,
	}

	log.Println("Opening serial port connection")
	port, error := serial.OpenPort(&config)

	if error != nil {
		log.Fatalln("Unable to open serial port connection: " + error.Error())
	}

	log.Println("Successfully connected to serial port")
	instance.port = port
	instance.readBuffer = make([]byte, serialPortReadBufferSizeBytes)
}

// Disconnect closes the active serial port
func (instance *Instance) Disconnect() {
	instance.port.Close()
}

// Show stored diagnostic trobule codes (DTC)
//
// References:
//
// - https://en.wikipedia.org/wiki/ISO_15765-2
//
// - https://en.wikipedia.org/wiki/OBD-II_PIDs#Service_03_-_Show_stored_Diagnostic_Trouble_Codes_(DTCs)
func (instance *Instance) ShowStoredDiagnosticTroubleCodes() []string {
	troubleCodes := make([]string, 0)
	data, error := instance.sendAndWaitForAnswer([]byte(modeShowStoredDtc))

	if error != nil {
		log.Println("Failed to fetch DTCs: " + error.Error())
		return troubleCodes
	}

	numberOfBytesReceived := len(data)

	if numberOfBytesReceived <= 4 {
		// Two or fewer DTCs received - process in single-frame mode
		for i := 0; i < len(data); i += 2 {
			troubleCodes = append(troubleCodes, BytesToDtc(data[i], data[i+1]))
		}
	} else {
		log.Println("Received more than 2 DTCs - this is not supported (yet)")
		log.Print("Raw binary data (decimal): ")
		log.Println(data)
	}

	return troubleCodes
}

// Send an OBD-II command to the ECU and wait for a response
//
// This operation will block until data is received, or the operation times out
func (instance *Instance) sendAndWaitForAnswer(command []byte) ([]byte, error) {
	_, error := instance.port.Write(command)

	// Unable to write to port
	if error != nil {
		log.Println("Failed to write to serial port: " + error.Error())
		return nil, error
	}

	size, _ := instance.port.Read(instance.readBuffer)

	// Read timeout
	if size == 0 {
		return nil, utility.GenericErrorMessage{Message: "Serial read timed out because zero bytes were returned"}
	}

	// Successful read
	log.Println("Successfully read " + strconv.Itoa(size) + " bytes from serial port")
	return instance.readBuffer[:size], nil
}
