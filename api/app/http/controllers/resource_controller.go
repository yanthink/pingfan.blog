package controllers

import (
	r "blog/app/http/requests/resources"
	"blog/app/http/responses"
	"blog/app/services"
	"github.com/gin-gonic/gin"
)

type resourceController struct {
}

func (*resourceController) Upload(c *gin.Context) {
	req := r.UploadValidate(c)

	responses.Json(c, gin.H{"url": services.Resource.Upload(req.File, req.Type)})
}
