package api

import (
	"encoding/json"
	"net/http"
	"time"
)

type systemStatusApi struct {
	Endpoint string
}

func (api systemStatusApi) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		api.get(response, request)
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (api systemStatusApi) get(response http.ResponseWriter, request *http.Request) {
	data := struct {
		CpuTemperature float32
		Timestamp      uint64
	}{
		CpuTemperature: -1.0,
		Timestamp:      uint64(time.Now().UnixMilli()),
	}

	json, error := json.Marshal(data)

	if error != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Header().Add(contentTypeHeader, mimeTypeText)

		response.Write([]byte(error.Error()))
	} else {
		response.Header().Add(contentTypeHeader, mimeTypeJson)
		response.Write(json)
	}
}
