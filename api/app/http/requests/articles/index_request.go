package articles

import (
	"blog/app/pagination"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
	"time"
)

type IndexRequest struct {
	pagination.Paginator
	UserID    int64      `form:"userId" json:"userId"`
	TagID     int64      `form:"tagId" json:"tagId"`
	TagIDs    int64      `form:"tagIds" json:"tagIds"`
	Order     string     `form:"order" json:"order"`
	StartDate *time.Time `form:"startDate" json:"startDate" time_format:"2006-01-02"`
	EndDate   *time.Time `form:"endDate" json:"endDate" time_format:"2006-01-02"`
}

func IndexValidate(c *gin.Context) (r *IndexRequest) {
	r = &IndexRequest{}
	validation.Validate(c, r)

	return
}
