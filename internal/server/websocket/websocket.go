package websocket

import (
	"github.com/gorilla/websocket"
	"log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func GetUpgrader() websocket.Upgrader {
	return upgrader
}

func (s *Subscription) WritePump() {
	conn := s.Connection
	defer func() {
		_ = conn.WebSocket.Close()
	}()

	for {
		select {
		case msg, ok := <-conn.Send:
			if !ok {
				log.Printf("Couldn't read message, closing socket")
				_ = conn.WebSocket.WriteMessage(websocket.CloseMessage, []byte{})
			}
			if err := conn.WebSocket.WriteJSON(msg); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (s *Subscription) ReadPump() {
	conn := s.Connection
	defer func() {
		hub.Unregister <- s
		_ = conn.WebSocket.Close()
	}()

	for {
		rmsg := &ReceivedMessage{}
		err := conn.WebSocket.ReadJSON(rmsg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		username, ok := GetUsername(s.UserID, rmsg.RoomID)
		if len(rmsg.Message) == 0 || !ok {
			log.Printf("Something wrong with message")
			continue
		}

		msg := &Message{
			RoomID:   rmsg.RoomID,
			Username: username,
			Message:  rmsg.Message,
		}
		hub.Broadcast <- msg
	}
}
