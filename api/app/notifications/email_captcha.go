package notifications

import (
	"blog/app/models"
	"fmt"
)

type EmailCaptcha struct {
	notification
	Email string
	Code  string
}

func (c *EmailCaptcha) Users() (users models.Users) {
	users = models.Users{&models.User{Email: &c.Email}}

	return
}

func (c *EmailCaptcha) Subject() string {
	return "邮箱验证码"
}

func (c *EmailCaptcha) Message() string {
	return fmt.Sprintf("验证码：%s，5分钟内有效。请勿向任何人泄露。", c.Code)
}

func (*EmailCaptcha) ToDatabase() (*map[string]any, bool) {
	return nil, false
}

func (*EmailCaptcha) ToWebsocket() bool {
	return false
}
