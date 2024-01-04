package mini_program

import "blog/config"

type IClient interface {
	Code2Session(code string, anonymousCode string) Session
}

type Session struct {
	SessionKey      string `json:"-"`
	Unionid         string `json:"unionid"`
	Openid          string `json:"openid"`
	AnonymousOpenid string `json:"anonymousOpenid"`
}

var Wx *Wechat

func init() {
	Wx = &Wechat{
		Appid:  config.MiniProgram.Wechat.Appid,
		Secret: config.MiniProgram.Wechat.Secret,
	}
}
