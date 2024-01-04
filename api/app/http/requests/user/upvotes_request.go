package user

import (
	"blog/app/pagination"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
)

type UpvotesRequest struct {
	pagination.Paginator
	UserID int64 `form:"-" json:"-"`
}

func UpvotesValidate(c *gin.Context) (r *UpvotesRequest) {
	r = &UpvotesRequest{}
	validation.Validate(c, r)

	r.UserID = c.GetInt64("userId")

	return
}
