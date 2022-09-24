package obdii

import (
	"log"
	"strconv"
	"time"

	"github.com/tarm/serial"
	"github.com/tntmeijs/pitona/server/utility"
)

const (
	debugSerialPortName               = "OBD-II_DEBUG"
	serialPortName                    = "/dev/ttyUSB0"
	serialPortBaudRate                = 9_600
	serialPortReadTimeoutMilliseconds = 5_000
	serialPortReadBufferSizeBytes     = 8
)

// Represents an OBD-II serial connection
type Instance struct {
	port       *serial.Port
	readBuffer []byte
}

// Connect establishes a connection with a serial port
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

	initialFrame := isoTpFrame{}
	if error := initialFrame.parse(data...); error != nil {
		log.Println("Unable to parse frame: " + error.Error())
		return troubleCodes
	}

	troubleCodes = append(troubleCodes, "TODO")

	return troubleCodes
}

// Send an OBD-II command to the ECU and wait for a response
//
// This operation will block until data is received, or until the operation times out
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

// Send an OBD-II command to the ECU and wait for N bytes to be returned
//
// This operation will block until all data is received, or until the operation times out
func (instance *Instance) sendAndWaitForMultipleAnswers(command []byte, expectedBytes int) ([]byte, error) {
	received := make([]byte, 0)

	for {
		data, error := instance.sendAndWaitForAnswer(command)

		if error != nil {
			return nil, error
		}

		received = append(received, data...)

		if len(received) >= expectedBytes {
			return received, nil
		}
	}
}
