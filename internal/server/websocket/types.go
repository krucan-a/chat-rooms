package websocket

import "github.com/gorilla/websocket"

type Hub struct {
	Rooms      map[string]map[*Subscription]bool
	Broadcast  chan *Message
	Register   chan *Subscription
	Unregister chan *Subscription
}

type Subscription struct {
	Connection *Connection
	UserID     string
	Username   string
	RoomID     string
}

type ReceivedMessage struct {
	RoomID  string `json:"room-id"`
	Message string `json:"message"`
}

type Message struct {
	RoomID   string `json:"room-id"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

type Connection struct {
	WebSocket *websocket.Conn
	Send      chan *Message
}
