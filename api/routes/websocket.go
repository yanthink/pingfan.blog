package routes

import (
	"blog/app/http/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterWebsocketRouter(r *gin.Engine) {
	r.GET("ws", controllers.Websocket.Ws)
}
