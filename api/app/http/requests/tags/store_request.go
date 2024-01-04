package tags

import (
	"blog/app/models"
	"blog/app/services"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
)

type StoreRequest struct {
	Name string `form:"name" json:"name" binding:"required"`
	Sort int64  `form:"sort" json:"sort" binding:"gte=0"`
}

func (r *StoreRequest) Store() *models.Tag {
	return services.Tag.Add(&models.Tag{
		Name: r.Name,
		Sort: r.Sort,
	})
}

func StoreValidate(c *gin.Context) (r *StoreRequest) {
	r = &StoreRequest{}
	validation.Validate(c, r)

	return
}
