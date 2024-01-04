package user

import (
	"blog/app/pagination"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
)

type FavoritesRequest struct {
	pagination.Paginator
	UserID int64 `form:"-" json:"-"`
}

func FavoritesValidate(c *gin.Context) (r *FavoritesRequest) {
	r = &FavoritesRequest{}
	validation.Validate(c, r)

	r.UserID = c.GetInt64("userId")

	return
}
