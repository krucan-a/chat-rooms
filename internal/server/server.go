package server

import (
	"chat-rooms/internal/server/routes"
	"github.com/gin-gonic/gin"
)

func registerRoutes(router *gin.Engine) *gin.Engine {
	routes.RegisterStatic(router)
	routes.RegisterAPI(router)
	routes.RegisterWebSocket(router)
	routes.RegisterRoom(router)
	return router
}

func CreateServer() *gin.Engine {
	return registerRoutes(gin.Default())
}
