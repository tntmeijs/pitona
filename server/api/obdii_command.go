package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tntmeijs/pitona/server/obdii"
)

type obdiiDiagnosticTroubleCodeApi struct {
	Endpoint      string
	ObdiiInstance *obdii.Instance
}

func (api obdiiDiagnosticTroubleCodeApi) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		api.get(response, request)
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (api obdiiDiagnosticTroubleCodeApi) get(response http.ResponseWriter, request *http.Request) {
	log.Println("Received request to fetch all stored DTCs")

	data := struct {
		Codes []string
	}{
		Codes: api.ObdiiInstance.ShowStoredDiagnosticTroubleCodes(),
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
