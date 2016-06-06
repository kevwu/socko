package socko

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
)

// optional, for convenience. serves the JS file to a client.
func ServeJS() {

}

// TODO: make clear the difference between a "message" and an "event".
// Update code to use this terminology correctly
type SockHandler struct {
	events  map[string]SockHandlerFunc
	onOpen  func()
	onClose func()
}

// data that an event handler can get ahold of
type SockEventData struct {
	messageType    string                 // type of event
	messageContent map[string]interface{} // JSON data sent with the message
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
func (sh *SockHandler) OnOpen(fn func()) {
	sh.onOpen = fn
}

func (sh *SockHandler) OnClose(fn func()) {
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

	// _, msgText, err := conn.ReadMessage()
	// // TODO: handle err
	// if err != nil {
	// 	break
	// }

	// msg, err := readMessageData(msgText)
	// if err != nil {
	// 	panic(err)
	// 	return
	// }

	sh.onOpen()

	for {
		_, msgText, err := conn.ReadMessage() // message type is expected to be TextMessage
		if err != nil {
			if err == websocket.ErrCloseSent {
				break
			}
		}

		msg, err := readMessageData(msgText)
		if err != nil {
			// TODO: log/report error
			panic(err)
			continue
		}

		messageType := msg.messageType

		if _, eventExists := sh.events[messageType]; eventExists {
			sh.events[messageType](msg)
		}
	}

	sh.onClose()
}

// not the best function name, consider renaming
// takes in the raw message string and returns a socket event data struct
// with the message content and type fields set.
func readMessageData(msgText []byte) (SockEventData, error) {
	msg := map[string]string{}
	err := json.Unmarshal(msgText, &msg)
	if err != nil {
		var empty SockEventData
		return empty, errors.New("Malformed JSON: " + err.Error())
	}

	if _, keySet := msg["type"]; !keySet {
		// received data doesn't have the event set
		var empty SockEventData
		return empty, errors.New("Received data is missing message type.")
	}

	var messageContent map[string]interface{}
	if _, keySet := msg["content"]; !keySet {
		var empty SockEventData
		return empty, errors.New("Received data is missing message content.")
	}

	err = json.Unmarshal([]byte(msg["content"]), &messageContent)

	return SockEventData{
		messageType:    msg["type"],
		messageContent: messageContent,
	}, nil
}
