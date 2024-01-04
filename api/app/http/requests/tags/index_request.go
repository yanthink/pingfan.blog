package tags

import (
	"blog/app/pagination"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
)

type IndexRequest struct {
	pagination.Paginator
	Q string `form:"q" json:"q"`
}

func IndexValidate(c *gin.Context) (r *IndexRequest) {
	r = &IndexRequest{}
	validation.Validate(c, r)

	return
}
