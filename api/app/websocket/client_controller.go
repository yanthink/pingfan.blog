package websocket

import (
	"blog/app"
	"blog/config"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func PingController(client *Client) {
	app.Logger.Debug("ping", zap.Uint64("UserID", client.UserID), zap.String("TempUserID", client.TempUserID))

	now := uint64(time.Now().Unix())
	client.Heartbeat(now)
}

func LoginController(client *Client, req *Request) {
	now := uint64(time.Now().Unix())
	token, ok := req.Data["token"].(string)
	if !ok {
		panic("缺少 token")
	}

	t, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(_ *jwt.Token) (any, error) {
		return config.Jwt.Key, nil
	})

	if err != nil {
		panic("token 校验失败")
	}

	var userId uint64
	var tempUserId string

	if claims, ok := t.Claims.(*jwt.RegisteredClaims); ok && t.Valid {
		if userId, _ = strconv.ParseUint(claims.Subject, 10, 64); userId == 0 {
			tempUserId = claims.Subject
		}
	}

	if userId == 0 && tempUserId == "" {
		panic("用户不存在")
	}

	if userId == client.UserID && tempUserId == client.TempUserID {
		panic("请勿重复登录")
	}

	clientManager.DelUser(client)
	clientManager.DelTempUser(client)

	if userId > 0 {
		client.UserID = userId
	}

	if tempUserId != "" {
		client.TempUserID = tempUserId
	}

	client.LoginTime = now
	client.Heartbeat(now)

	clientManager.Login <- client

	client.SendMsg(&Response{Event: Login, Data: "success"})
}

func LogoutController(client *Client) {
	clientManager.DelUser(client)
	clientManager.DelTempUser(client)

	app.Logger.Debug("用户退出", zap.Uint64("UserID", client.UserID), zap.String("TempUserID", client.TempUserID))

	client.SendMsg(&Response{Event: Logout, Data: "success"})
}
