package articles

import (
	"blog/app/helpers"
	h "blog/app/http"
	"blog/app/http/responses"
	"blog/app/models"
	"blog/app/services"
	"blog/app/validation"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UpdateRequest struct {
	StoreRequest
	Original *models.Article
}

func (r *UpdateRequest) Update(id int64) *models.Article {
	article := models.Article{
		Title:   r.Title,
		Content: r.Content,
		Preview: r.Preview,
		Tags: helpers.Map(r.TagIds, func(_ int, id json.Number) *models.Tag {
			tagId, err := id.Int64()
			if err != nil {
				panic("标签不存在！")
			}

			return &models.Tag{ID: tagId}
		}),
	}

	return services.Article.Update(id, &article, r.Original)
}

func UpdateValidate(c *gin.Context, id int64) (r *UpdateRequest) {
	r = &UpdateRequest{}
	validation.Validate(c, r)

	article := services.Article.GetByID(id)

	if c.GetInt64("userId") != article.UserID {
		panic(&h.Error{StatusCode: http.StatusForbidden, Code: responses.CodeAccessDenied, Message: "无权限操作"})
	}

	r.Original = article

	return
}
