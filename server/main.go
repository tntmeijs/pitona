package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const ContentTypeHeader = "Content-Type"
const MimeTypeText = "text/plain"
const MimeTypeJson = "application/json"

type systemStatusResponse struct {
	CpuTemperature float32
	Timestamp      uint64
}

type serialListAvailablePortsResponse struct {
	Ports []string
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

func handleSerialConnectionRequest(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		getHandleSerialConnectionRequest(response)
	case http.MethodPost:
		postHandleSerialConnectionRequest(response)
	case http.MethodDelete:
		deleteHandleSerialConnectionRequest(response)
	}
}

func getHandleSerialConnectionRequest(response http.ResponseWriter) {
	data := serialListAvailablePortsResponse{Ports: []string{"PORT NAME HERE"}}

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

func postHandleSerialConnectionRequest(response http.ResponseWriter) {
	log.Println("POST")
}

func deleteHandleSerialConnectionRequest(response http.ResponseWriter) {
	log.Println("DELETE")
}

func handleObdiiCommandRequest(response http.ResponseWriter, request *http.Request) {
	log.Println("COMMAND")
}

func main() {
	http.HandleFunc("/api/v1/system/status", handleSystemStatusRequest)
	http.HandleFunc("/api/v1/serial/port", handleSerialConnectionRequest)
	http.HandleFunc("/api/v1/obdii/command", handleObdiiCommandRequest)

	log.Println("Starting HTTP server on port 80")
	log.Fatal(http.ListenAndServe("localhost:80", nil))
}
