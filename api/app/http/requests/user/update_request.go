package user

import (
	"blog/app/services"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
)

type UpdateRequest struct {
	Status *int64 `form:"status" json:"status"`
}

func UpdateValidate(c *gin.Context) (r *UpdateRequest) {
	services.User.CheckAuthIsAdmin(c)

	r = &UpdateRequest{}
	validation.Validate(c, r)

	return
}
