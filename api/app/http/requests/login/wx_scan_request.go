package login

import (
	"blog/app/models"
	"blog/app/services"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
)

type WxScanRequest struct {
	Uuid string       `form:"uuid" json:"uuid" binding:"required"`
	User *models.User `form:"-" json:"-"`
}

func (r *WxScanRequest) Login() {
	services.User.WxScanLogin(r.User, r.Uuid)
}

func WxScanValidate(c *gin.Context) (r *WxScanRequest) {
	r = &WxScanRequest{}
	validation.Validate(c, r)

	r.User = services.User.GetAuthUser(c)

	return
}
