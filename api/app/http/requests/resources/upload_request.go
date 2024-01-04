package resources

import (
	"blog/app/resource"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

type UploadRequest struct {
	Type resource.Type         `form:"type" binding:"required"`
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func UploadValidate(c *gin.Context) (r *UploadRequest) {
	r = &UploadRequest{}
	validation.Validate(c, r)

	r.Type.Validate(r.File)

	return
}
