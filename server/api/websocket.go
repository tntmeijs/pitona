package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tntmeijs/pitona/server/obdii"
)

const (
	// Time allowed to write a message to a peer
	WriteWait = 3 * time.Second

	// Maximum size of a message in bytes
	MaxMessageSize = 8_192

	// Time allowed to read the next pong message from the peer
	PongWait = 60 * time.Second

	// Send a ping to the peer once a period
	PingPeriod = (PongWait * 9) / 10
)

type webSocketHandler struct {
	ecuRequests  chan<- obdii.EcuRequest
	ecuResponses <-chan obdii.EcuResponse
	settings     websocket.Upgrader
}

func NewWebSocketHandler(ecuRequests chan<- obdii.EcuRequest, ecuResponses <-chan obdii.EcuResponse, settings websocket.Upgrader) webSocketHandler {
	return webSocketHandler{
		ecuRequests:  ecuRequests,
		ecuResponses: ecuResponses,
		settings:     settings,
	}
}

func (handler *webSocketHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	connection, err := handler.settings.Upgrade(writer, request, nil)

	if err != nil {
		log.Println("WebSocket connection upgrade failed", err)
		return
	}

	connection.SetReadLimit(MaxMessageSize)

	defer func() {
		connection.Close()
		log.Println("WebSocket connection with client has been closed")
	}()

	log.Println("WebSocket connection established with client")

	go webSocketWriter(connection, handler.ecuResponses)
	webSocketReader(connection, handler.ecuRequests)
}

func webSocketWriter(ws *websocket.Conn, ecuResponses <-chan obdii.EcuResponse) {
	pingTicker := time.NewTicker(PingPeriod)
	defer func() { pingTicker.Stop() }()

	for {
		select {
		case ecuResponse := <-ecuResponses:
			ws.SetWriteDeadline(time.Now().Add(WriteWait))

			if err := ws.WriteJSON(ecuResponse); err != nil {
				// Error writing data, stop the writer to force the client to reconnect
				log.Println("Failed to write ECU response:", err)
				return
			}
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(WriteWait))

			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				// Could not receive pong from client, time to close the writer
				return
			}
		}
	}
}

func webSocketReader(ws *websocket.Conn, ecuRequests chan<- obdii.EcuRequest) {
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(PongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(PongWait)); return nil })

	ws.SetReadDeadline(time.Now().Add(PongWait))
	ws.SetPongHandler(func(appData string) error {
		ws.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})

	for {
		var ecuRequest obdii.EcuRequest
		err := ws.ReadJSON(&ecuRequest)

		if err != nil {
			if !websocket.IsCloseError(err) && !websocket.IsUnexpectedCloseError(err) {
				log.Println("Unexpected error in WebSocket read function:", err)
			}

			log.Println("Cannot read message:", err)
			// Connection was closed
			return
		}

		log.Println(ecuRequest)

		// Write the ECU request message to the ECU communication channel to ensure the ECU will
		// pick it up next time it fetches all pending requests
		ecuRequests <- ecuRequest
	}
}
