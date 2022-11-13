package obdii

type PidRequest struct {
	Mode                string
	Pid                 string
	ResponseSizeInBytes int
}

type PidResponse struct {
	Mode string
	Pid  string
	Data string
}
