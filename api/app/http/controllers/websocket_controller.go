package controllers

import (
	"blog/app/websocket"
	"github.com/gin-gonic/gin"
)

type websocketController struct {
}

func (*websocketController) Ws(c *gin.Context) {
	websocket.Upgrade(c.Writer, c.Request)
}
