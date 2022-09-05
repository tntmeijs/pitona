package obdii

const modeShowStoredDtc = "03"

// Represents an OBD-II PID
//
// Reference: https://en.wikipedia.org/wiki/OBD-II_PIDs
type Pid interface {
	mode() string
	pid() string
	responseSizeInBytes() uint8
}
