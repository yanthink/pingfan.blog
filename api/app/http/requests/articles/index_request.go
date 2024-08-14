package articles

import (
	"blog/app/models"
	"blog/app/pagination"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
	"time"
)

type IndexRequest struct {
	pagination.Paginator `filter:"-"`
	UserID               int64       `form:"userId"`
	Order                string      `form:"order"`
	StartDate            *time.Time  `form:"startDate" time_format:"2006-01-02"`
	EndDate              *time.Time  `form:"endDate" time_format:"2006-01-02"`
	Tag                  *models.Tag `form:"tag"`
}

func IndexValidate(c *gin.Context) (r *IndexRequest) {
	r = &IndexRequest{}
	validation.Validate(c, r)

	return
}
