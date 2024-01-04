package user

import (
	"blog/app/pagination"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
)

type LikesRequest struct {
	pagination.Paginator
	UserID int64 `form:"-" json:"-"`
}

func LikesValidate(c *gin.Context) (r *LikesRequest) {
	r = &LikesRequest{}
	validation.Validate(c, r)

	r.UserID = c.GetInt64("userId")

	return
}
