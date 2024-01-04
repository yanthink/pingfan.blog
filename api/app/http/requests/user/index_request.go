package user

import (
	"blog/app/models"
	"blog/app/pagination"
	"blog/app/services"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
)

type IndexRequest struct {
	pagination.Paginator
	ID     int64            `form:"id" json:"id"`
	Name   string           `form:"name" json:"name"`
	Email  string           `form:"email" json:"email"`
	Openid string           `form:"openid" json:"openid"`
	Role   *models.UserRole `form:"role" json:"role"`
	Status *int64           `form:"status" json:"status"`
}

func IndexValidate(c *gin.Context) (r *IndexRequest) {
	services.User.CheckAuthIsAdmin(c)

	r = &IndexRequest{}
	validation.Validate(c, r)

	return
}
