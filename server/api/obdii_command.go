package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/tntmeijs/pitona/server/obdii"
)

// Simplifies debugging of the interactions with the ECU
type obdiiDebugApi struct {
	ObdiiInstance *obdii.Instance
}

// Allows interactions of type service 01
type odbiiMode01Api struct {
	ObdiiInstance *obdii.Instance
}

// Allows interactions of type service 03
type obdiiMode03Api struct {
	ObdiiInstance *obdii.Instance
}

func (api obdiiDebugApi) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	log.Println("Received request to send raw data to ECU (debug mode)")

	defer request.Body.Close()
	body, error := io.ReadAll(request.Body)

	if error != nil {
		log.Println("Unable to read request's body: " + error.Error())
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(error.Error()))
		return
	}

	requestBody := struct {
		Command string
	}{}

	if json.Unmarshal(body, &requestBody) != nil {
		log.Println("Unable to unmarshal request body: " + error.Error())
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(error.Error()))
		return
	}

	log.Printf("Sending command \"%s\" to ECU\n", requestBody.Command)
	data, error := api.ObdiiInstance.SendAndWaitForAnswer([]byte(requestBody.Command), -1)

	if error != nil {
		log.Println("Unable to fetch data from ECU: " + error.Error())
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(error.Error()))
		return
	}

	if len(data) == 0 {
		log.Println("No data received from ECU")
		response.WriteHeader(http.StatusServiceUnavailable)
		response.Write([]byte("No data received from ECU"))
		return
	}

	rawDecStr := []string{}
	rawHexStr := []string{}
	rawBinStr := []string{}

	// Discard index here - without using an index, Go will implicitly convert each item in data
	// into integers
	for _, value := range data {
		rawDecStr = append(rawDecStr, strconv.Itoa(int(value)))
		rawHexStr = append(rawHexStr, hex.EncodeToString([]byte{byte(value)}))
		rawBinStr = append(rawBinStr, fmt.Sprintf("%b", value))
	}

	log.Printf("Raw DEC: %v\n", rawDecStr)
	log.Printf("Raw HEX: %v\n", rawHexStr)
	log.Printf("Raw BIN: %v\n", rawBinStr)

	responseBody := struct {
		Length             int
		Command            string
		CommandTranslation string
		RawDec             []string
		RawHex             []string
		RawBin             []string
	}{
		Length:             len(data),
		Command:            requestBody.Command,
		CommandTranslation: "TODO: translate PID into English",
		RawDec:             rawDecStr,
		RawHex:             rawHexStr,
		RawBin:             rawBinStr,
	}

	responseData, error := json.Marshal(responseBody)

	if error != nil {
		log.Println("Unable to marshal data from ECU: " + error.Error())
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(error.Error()))
		return
	}

	response.Header().Add(contentTypeHeader, mimeTypeJson)
	response.Write(responseData)
}

func (api odbiiMode01Api) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	log.Println("Received request to show current data (service 01)")
}

func (api obdiiMode03Api) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	log.Println("Received request to fetch all stored DTCs (service 03)")

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
