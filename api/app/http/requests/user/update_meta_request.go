package user

import (
	"blog/app/validation"
	"github.com/gin-gonic/gin"
)

type UpdateMetaRequest struct {
	Meta Meta `form:"meta" json:"meta" binding:"required"`
}

type Meta struct {
	EmailNotify int64 `form:"emailNotify" json:"emailNotify" binding:"oneof=0 1 2"`
}

func UpdateMetaValidate(c *gin.Context) (r *UpdateMetaRequest) {
	r = &UpdateMetaRequest{}
	validation.Validate(c, r)

	return
}
