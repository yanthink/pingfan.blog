package websocket

import (
	"blog/app"
	ws "github.com/gorilla/websocket"
	"net/http"
	"time"
)

var upgrade = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clientManager = NewClientManager()

func Start() {
	go clientManager.start()
}

func Upgrade(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	now := uint64(time.Now().Unix())
	client := NewClient(conn.RemoteAddr().String(), conn, now)

	// 用户连接事件
	clientManager.Register <- client
	go func() {
		// 监听连接“完成”事件，其实也可以说丢失事件
		<-r.Context().Done()
		app.Logger.Sugar().Debug("socket 断开")
		clientManager.Unregister <- client
	}()

	defer func() {
		if r := recover(); r != nil {
			app.Logger.Sugar().Errorf("read stop: %v", r)
		}
	}()

	for {
		_, message, err := client.Socket.ReadMessage()
		if err != nil {
			return
		}
		client.MessageHandle(message)
	}
}
