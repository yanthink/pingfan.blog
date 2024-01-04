package articles

import (
	"blog/app/helpers"
	"blog/app/models"
	"blog/app/services"
	"blog/app/validation"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type StoreRequest struct {
	Title   string        `form:"title" json:"title" label:"标题" binding:"required"`
	Content string        `form:"title" json:"content" label:"内容" binding:"required"`
	Preview string        `form:"preview" json:"preview" label:"预览图" binding:"omitempty,url"`
	TagIds  []json.Number `form:"tagIds" json:"tagIds" label:"标签" binding:"omitempty,gt=0,lte=10,dive,min=1"`
	User    *models.User  `form:"-" json:"-"`
}

func (r *StoreRequest) Store() *models.Article {
	article := models.Article{
		UserID:  r.User.ID,
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

	return services.Article.Add(&article)
}

func StoreValidate(c *gin.Context) (r *StoreRequest) {
	user := services.User.CheckAuthIsAdmin(c)

	r = &StoreRequest{}
	validation.Validate(c, r)

	r.User = user

	return
}
