package obdii

// Response data from the ECU
type EcuResponse struct {
	Pids []PidResponse
}

// Payload sent to the ECU from a client
type EcuRequest struct {
	Pids []PidRequest
}
