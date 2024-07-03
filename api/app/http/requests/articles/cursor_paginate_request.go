package articles

import (
	"blog/app/pagination"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
	"time"
)

type CursorPaginateRequest struct {
	pagination.CursorPaginator
	UserID    int64      `form:"userId" json:"userId"`
	TagID     int64      `form:"tagId" json:"tagId"`
	TagIDs    string     `form:"tagIds" json:"tagIds"`
	Order     string     `form:"order" json:"order"`
	StartDate *time.Time `form:"startDate" json:"startDate" time_format:"2006-01-02"`
	EndDate   *time.Time `form:"endDate" json:"endDate" time_format:"2006-01-02"`
}

func CursorPaginateValidate(c *gin.Context) (req *CursorPaginateRequest) {
	req = &CursorPaginateRequest{}
	validation.Validate(c, req)

	rules := []paginator.Rule{
		{
			Key:     "Hotness",
			Order:   paginator.DESC,
			SQLRepr: "articles.hotness",
		},
		{
			Key:     "ID",
			Order:   paginator.DESC,
			SQLRepr: "articles.id",
		},
	}

	switch req.Order {
	case "latest":
		rules = []paginator.Rule{
			{
				Key:     "ID",
				Order:   paginator.DESC,
				SQLRepr: "articles.id",
			},
		}
	case "like":
		rules = []paginator.Rule{
			{
				Key:     "LikeCount",
				Order:   paginator.DESC,
				SQLRepr: "articles.like_count",
			},
		}
	case "comment":
		rules = []paginator.Rule{
			{
				Key:     "CommentCount",
				Order:   paginator.DESC,
				SQLRepr: "articles.comment_count",
			},
		}
	}

	req.CursorPaginator.Config.Rules = rules

	return
}
