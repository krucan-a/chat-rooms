package main

import (
	"chat-rooms/internal/server"
	"chat-rooms/internal/server/websocket"
	"os"
)

func main() {
	go websocket.GetHub().Run()

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = ":8080"
	}
	_ = server.CreateServer().Run(port)
}
