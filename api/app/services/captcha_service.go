package services

import (
	"blog/app/captcha"
	"blog/app/mail"
	"fmt"
)

type captchaService struct {
}

func (*captchaService) SendEmailCaptcha(email string) string {
	code, hash := captcha.New().Generate(email)

	err := mail.New().
		AppendTo(email).
		Send("邮箱验证码", fmt.Sprintf("验证码：%s，5分钟内有效。请勿向任何人泄露。", code))

	if err != nil {
		panic(err)
	}

	return hash
}
