package main

import (
	"chat-rooms/internal/server"
	"chat-rooms/internal/server/websocket"
)

func main() {
	go websocket.GetHub().Run()

	_ = server.CreateServer().Run()
}
