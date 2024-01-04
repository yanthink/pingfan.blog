package tags

import (
	"blog/app/models"
	"blog/app/services"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
)

type UpdateRequest struct {
	Name string `form:"name" json:"name" binding:"required"`
	Sort int64  `form:"sort" json:"sort" binding:"gte=0"`
}

func (r *UpdateRequest) Update(id int64) *models.Tag {
	return services.Tag.Update(id, &models.Tag{
		Name: r.Name,
		Sort: r.Sort,
	})
}

func UpdateValidate(c *gin.Context) (r *UpdateRequest) {
	r = &UpdateRequest{}
	validation.Validate(c, r)

	return
}
