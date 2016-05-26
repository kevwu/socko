package socko

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

// optional, for convenience. serves the JS file to a client.
func ServeJS() {

}

type SockHandler struct {
	events  map[string]SockHandlerFunc
	onOpen  SockHandlerFunc
	onClose SockHandlerFunc
}

// data that an event handler can get ahold of
type SockEventData struct {
	eventType string                 // type of event
	msgData   map[string]interface{} // JSON data sent with the message
}

// represents an individual connection
type SockConnection struct {
	ws   *websocket.Conn
	uuid string
}

type SockHandlerFunc func(SockEventData)

// initializes a new socket handler
func NewSockHandler() *SockHandler {
	sh := &SockHandler{
		events: make(map[string]SockHandlerFunc),
	}
	return sh
}

// analogous to net/http's HandleFunc()
func (sh *SockHandler) OnEvent(event string, fn SockHandlerFunc) {
	sh.events[event] = fn
}

// runs when socket connection is open
func (sh *SockHandler) OnOpen(fn SockHandlerFunc) {
	sh.onOpen = fn
}

func (sh *SockHandler) OnClose(fn SockHandlerFunc) {
	sh.onClose = fn
}

// this is the HTTP handler function that gets passed into the router.
func (sh *SockHandler) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic("socko: Failed to upgrade connection to websocket. " + err.Error())
	}

	sh.onOpen()

	for {
		_, p, err = conn.ReadMessage() // message type is expected to be TextMessage
		if err != nil {
			if err == websocket.ErrCloseSent {
				break
			}
		}

		msg := map[string]string{}
		err = json.Unmarshal(p, &msg)
		if err != nil {
			panic("socko: Malformed JSON. " + err.Error())
			continue
		}

		var event string

		if event, keySet := msg["event"]; !keySet {
			// received data doesn't have the event set
			panic("socko: Received data is missing event value.")
		}

		if _, eventExists := sh.events[event]; eventExists {
			sh.events[event]()
		}
	}

	sh.onClose()
}
