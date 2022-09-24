package api

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/tntmeijs/pitona/server/obdii"
)

const contentTypeHeader = "Content-Type"
const mimeTypeText = "text/plain"
const mimeTypeJson = "application/json"

// ApiServer represents a simple REST API server that allows the outside world to send and receive
// commands from the motorcycle's onboard ECU
type ApiServer struct {
	ObdiiInstance    *obdii.Instance
	server           http.Server
	initiateShutdown chan struct{}
	shutdownComplete chan struct{}
}

// Start the API server and block until the server (gracefully) shuts down
//
// The server can be stopped by sending a DELETE request to /server
func (apiServer *ApiServer) Start(isDebug bool, port int) {
	apiServer.initiateShutdown = make(chan struct{})
	apiServer.shutdownComplete = make(chan struct{})

	log.Println("Launching API server")

	address := ""

	// For debugging purposes only
	if isDebug {
		address = "localhost"
	}

	endpoints := make(map[string]http.Handler)

	// Configure endpoint routing
	endpoints["/api/v1/obdii/debug"] = obdiiDebugApi{ObdiiInstance: apiServer.ObdiiInstance}
	endpoints["/api/v1/obdii/01"] = odbiiMode01Api{ObdiiInstance: apiServer.ObdiiInstance}
	endpoints["/api/v1/obdii/03"] = obdiiMode03Api{ObdiiInstance: apiServer.ObdiiInstance}
	endpoints["/api/v1/system/status"] = systemStatusApi{}

	// Register all endpoints with the webserver
	log.Println("Registering all HTTP handlers")
	for endpoint, handler := range endpoints {
		http.Handle(endpoint, handler)
	}

	// Register a special endpoint that can be used to gracefully stop the server
	http.HandleFunc("/server", func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodDelete {
			response.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		log.Println("Sending server shutdown signal")
		apiServer.RequestShutdown()
	})

	// Start a Goroutine that waits for a shutdown signal on the shutdown channel
	go apiServer.listenForShutdownSignal()

	// Start the webserver
	apiServer.server.Addr = address + ":" + strconv.Itoa(port)
	log.Println("Starting webserver (" + apiServer.server.Addr + ")")
	error := apiServer.server.ListenAndServe()

	if error != http.ErrServerClosed {
		log.Fatalln("Server closed unexpectedly: " + error.Error())
	}

	log.Println("Waiting for server shutdown")
	<-apiServer.shutdownComplete
	log.Println("Server shutdown complete")
}

// RequestShutdown sends a signal to the server to gracefully start shutting down
func (apiServer *ApiServer) RequestShutdown() {
	apiServer.initiateShutdown <- struct{}{}
}

// listenForShutdownSignal is a utility method that blocks until it receives a shutdown signal
//
// This method should always run as a Goroutine to prevent blocking the main thread
func (apiServer *ApiServer) listenForShutdownSignal() {
	<-apiServer.initiateShutdown
	log.Println("Shutdown initiation signal received")

	if error := apiServer.server.Shutdown(context.Background()); error != nil {
		log.Println("Failed to gracefully stop API server: " + error.Error())
	}

	// The application can safely exit once this channel has been signaled
	close(apiServer.shutdownComplete)
}
