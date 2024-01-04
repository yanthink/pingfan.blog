package user

import (
	"blog/app/pagination"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
)

type CommentsRequest struct {
	pagination.Paginator
	Type   *int64 `form:"type" json:"type" binding:"oneof=0 1"`
	UserID int64  `form:"-" json:"-"`
	Sort   string `form:"-" json:"-"`
}

func CommentsValidate(c *gin.Context) (r *CommentsRequest) {
	r = &CommentsRequest{}
	validation.Validate(c, r)

	r.UserID = c.GetInt64("userId")
	r.Sort = "DESC"

	return
}
