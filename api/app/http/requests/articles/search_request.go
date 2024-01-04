package articles

import (
	"blog/app/pagination"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
)

type SearchRequest struct {
	pagination.Paginator
	Q           string   `form:"q" json:"q"`
	Filters     []string `form:"filters" json:"filters"`
	QueryFields []string `form:"queryFields" json:"queryFields"`
	SortFields  []string `form:"sortFields" json:"sortFields"`
}

func SearchValidate(c *gin.Context) (r *SearchRequest) {
	r = &SearchRequest{}
	validation.Validate(c, r)

	if r.Current != nil {
		r.Page = *r.Current
	}

	if r.PageSize != nil {
		r.Limit = *r.PageSize
	}

	return
}
