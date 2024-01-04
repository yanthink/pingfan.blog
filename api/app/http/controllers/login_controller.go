package controllers

import (
	r "blog/app/http/requests/login"
	"blog/app/http/responses"
	"blog/app/models"
	"blog/app/services"
	"github.com/gin-gonic/gin"
	"time"
)

type loginController struct {
}

type LoginResponse struct {
	*models.User
	Token string `json:"token"`
}

func (*loginController) Account(c *gin.Context) {
	token, user := r.AccountValidate(c).Login()
	responses.Json(c, &LoginResponse{
		User:  user,
		Token: token,
	})
}

func (*loginController) Wx(c *gin.Context) {
	token, user := r.WxValidate(c).Login()
	responses.Json(c, &LoginResponse{
		User:  user,
		Token: token,
	})
}

func (*loginController) WxQRCode(c *gin.Context) {
	token, img, expiration := services.User.GetWxLoginQRCode()
	responses.Json(c, gin.H{"token": token, "img": img, "expiresIn": int64(expiration / time.Millisecond)})
}

func (*loginController) WxScan(c *gin.Context) {
	r.WxScanValidate(c).Login()
	responses.Json(c)
}
