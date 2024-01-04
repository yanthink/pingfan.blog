package login

import (
	"blog/app/models"
	"blog/app/services"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
)

type WxRequest struct {
	Code string `form:"code" json:"code" binding:"required"`
}

func (r *WxRequest) Login() (string, *models.User) {
	return services.User.WxLogin(r.Code)
}

func WxValidate(c *gin.Context) (r *WxRequest) {
	r = &WxRequest{}
	validation.Validate(c, r)

	return
}
