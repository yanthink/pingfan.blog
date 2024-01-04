package user

import (
	"blog/app/pagination"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
)

type NotificationsRequest struct {
	pagination.Paginator
	UserID int64 `form:"-" json:"-"`
}

func NotificationsValidate(c *gin.Context) (r *NotificationsRequest) {
	r = &NotificationsRequest{}
	validation.Validate(c, r)

	r.UserID = c.GetInt64("userId")

	return
}
