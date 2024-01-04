package websocket

import (
	"blog/app"
	ws "github.com/gorilla/websocket"
)

// Client 用户连接
type Client struct {
	Addr          string   // 客户端地址
	Socket        *ws.Conn // 用户连接
	UserID        uint64   // 用户ID，用户登录以后才有
	TempUserID    string   // 临时用户ID
	FirstTime     uint64   // 首次连接事件
	HeartbeatTime uint64   // 用户上次心跳时间
	LoginTime     uint64   // 登录时间，登录以后才有
}

func NewClient(addr string, socket *ws.Conn, firstTime uint64) *Client {
	return &Client{
		Addr:          addr,
		Socket:        socket,
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
	}
}

// SendMsg 发送消息
func (c *Client) SendMsg(msg *Response) {
	if c == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			app.Logger.Sugar().Errorf("SendMsg stop: %v", r)
		}
	}()

	_ = c.Socket.WriteJSON(msg)
}

// Heartbeat 用户心跳
func (c *Client) Heartbeat(t uint64) {
	c.HeartbeatTime = t
	return
}
