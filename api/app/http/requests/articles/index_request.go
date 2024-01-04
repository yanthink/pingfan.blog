package articles

import (
	h "blog/app/http"
	"blog/app/http/responses"
	"blog/app/pagination"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type IndexRequest struct {
	pagination.Paginator
	UserID    int64  `form:"userId" json:"userId"`
	TagID     int64  `form:"tagId" json:"tagId"`
	TagIDs    int64  `form:"tagIds" json:"tagIds"`
	Order     string `form:"order" json:"order"`
	StartDate string `form:"startDate" json:"startDate"`
	EndDate   string `form:"endDate" json:"endDate"`
}

func IndexValidate(c *gin.Context) (r *IndexRequest) {
	r = &IndexRequest{}
	validation.Validate(c, r)

	if r.StartDate != "" {
		if _, err := time.Parse(time.DateOnly, r.StartDate); err != nil {
			panic(&h.Error{Code: responses.CodeParameterError, StatusCode: http.StatusUnprocessableEntity, Message: "时间格式不正确"})
		}
	}

	if r.EndDate != "" {
		if _, err := time.Parse(time.DateOnly, r.EndDate); err != nil {
			panic(&h.Error{Code: responses.CodeParameterError, StatusCode: http.StatusUnprocessableEntity, Message: "时间格式不正确"})
		}
	}

	return
}
