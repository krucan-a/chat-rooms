package routes

import (
	wsocket "chat-rooms/internal/server/websocket"
	"chat-rooms/internal/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type APIMessage struct {
	Message string `json:"message"`
}

func RegisterAPI(router *gin.Engine) {
	api := router.Group("/api")
	{
		health(api)
		createRoom(api)
		checkRoom(api)
		registerUser(api)
		checkCookie(api)
	}
}

func health(apiRoute *gin.RouterGroup) {
	apiRoute.GET("/health", func(c *gin.Context) {
		c.JSON(200, "{\"status\":\"ok\"}")
	})
}

func createRoom(apiRoute *gin.RouterGroup) {
	apiRoute.GET("/create-room", func(c *gin.Context) {
		hub := wsocket.GetHub()

		newRoomID := util.GenerateID()
		for hub.Rooms[newRoomID] != nil {
			newRoomID = util.GenerateID()
		}

		hub.Rooms[newRoomID] = make(map[*wsocket.Subscription]bool)

		log.Printf("Created room %s", newRoomID)

		c.JSON(http.StatusOK, struct {
			RoomID string `json:"room-id"`
		}{newRoomID})
	})
}

func checkRoom(apiRoute *gin.RouterGroup) {
	apiRoute.GET("/check-room", func(c *gin.Context) {
		id := c.Query("id")

		if !wsocket.CheckRoomExists(id) {
			log.Printf("Room %s doesn't exist", id)
			c.JSON(http.StatusNotFound, APIMessage{"Room does not exist"})
		} else {
			c.JSON(http.StatusOK, nil)
		}
	})
}

func registerUser(apiRoute *gin.RouterGroup) {
	apiRoute.GET("/register-user", func(c *gin.Context) {
		username := c.Query("username")
		roomID := c.Query("room-id")

		if !wsocket.CheckUsernameAvailable(username, roomID) {
			log.Printf("User %s already exists in room %s", username, roomID)
			c.JSON(http.StatusConflict, APIMessage{"This username is already taken"})
		}

		cookieName := util.GenerateCookieName(roomID)
		userID := util.GenerateID()
		log.Printf("Registering user %s as %s", userID, username)
		c.SetCookie(cookieName, userID, 0, "/", "localhost", false, true)
	})
}

func checkCookie(apiRoute *gin.RouterGroup) {
	apiRoute.GET("/check-cookie", func(c *gin.Context) {
		roomID := c.Query("room-id")
		cookieName := util.GenerateCookieName(roomID)

		userID, err := c.Cookie(cookieName)
		username, ok := wsocket.GetUsername(userID, roomID)
		if err != nil || !ok {
			log.Printf("Cookie %s or user %s (%s) doesn't exist", cookieName, userID, username)
			c.JSON(http.StatusForbidden, nil)
		} else {
			c.JSON(http.StatusOK, APIMessage{username})
		}
	})
}
