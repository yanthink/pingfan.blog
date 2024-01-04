package websocket

import (
	"blog/app"
	"encoding/json"
	"fmt"
)

const (
	Ping   = "ping"
	Login  = "login"
	Logout = "logout"
)

func (c *Client) MessageHandle(message []byte) {
	req := &Request{}

	if err := json.Unmarshal(message, req); err != nil {
		app.Logger.Sugar().Errorf("请求参数错误：%v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			c.SendMsg(&Response{Event: req.Event, Data: fmt.Sprintf("%v", r)})
		}
	}()

	switch req.Event {
	case Ping:
		PingController(c)
		break
	case Login:
		LoginController(c, req)
		break
	case Logout:
		LogoutController(c)
		break
	}
}
