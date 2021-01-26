package routes

import (
	wsocket "chat-rooms/internal/server/websocket"
	"chat-rooms/internal/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func RegisterWebSocket(router *gin.Engine) {
	router.GET("/ws", func(c *gin.Context) {
		roomID := c.Query("room-id")
		username := c.Query("username")
		userID, err := c.Cookie(util.GenerateCookieName(roomID))
		if err != nil {
			log.Printf("User %s (%s) unknown", userID, username)
			c.JSON(http.StatusForbidden, APIMessage{"Couldn't verify user"})
			return
		}

		connectWebSocket(c.Writer, c.Request, userID, username, roomID)
	})
}

func connectWebSocket(w http.ResponseWriter, r *http.Request, userID, username, roomID string) {
	upgrader := wsocket.GetUpgrader()
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	conn := &wsocket.Connection{
		WebSocket: ws,
		Send:      make(chan *wsocket.Message),
	}
	sub := &wsocket.Subscription{
		Connection: conn,
		UserID:     userID,
		Username:   username,
		RoomID:     roomID,
	}

	wsocket.GetHub().Register <- sub

	go sub.WritePump()
	go sub.ReadPump()
}
