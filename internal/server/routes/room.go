package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterRoom(router *gin.Engine) {
	router.GET("/room/:id", func(context *gin.Context) {
		context.HTML(http.StatusOK, "room.html", nil)
	})
}
