package controllers

import (
	r "blog/app/http/requests/captcha"
	"blog/app/http/responses"
	"blog/app/services"
	"github.com/gin-gonic/gin"
)

type captchaController struct {
}

func (*captchaController) Email(c *gin.Context) {
	req := r.EmailValidate(c)

	responses.Json(c, gin.H{"key": services.Captcha.SendEmailCaptcha(req.Email)})
}
