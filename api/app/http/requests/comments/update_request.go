package comments

import (
	"blog/app/models"
	"blog/app/services"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
)

type UpdateRequest struct {
	Content string `form:"title" json:"content" label:"内容" binding:"required,max=512"`
}

func (r *UpdateRequest) Update(id int64) *models.Comment {
	return services.Comment.Update(id, &models.Comment{
		Content: r.Content,
	})
}

func UpdateValidate(c *gin.Context) (r *UpdateRequest) {
	services.User.CheckAuthIsAdmin(c)

	r = &UpdateRequest{}
	validation.Validate(c, r)

	return
}
