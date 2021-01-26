package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("./website/*.html")
	static := router.Group("/")
	{
		static.GET("/", func(context *gin.Context) {
			context.HTML(http.StatusOK, "index.html", nil)
		})

		static.StaticFS("/static", http.Dir("./website"))
	}
}
