package controllers

import (
	r "blog/app/http/requests/tags"
	"blog/app/http/responses"
	"blog/app/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type tagController struct {
}

func (*tagController) Index(c *gin.Context) {
	tags, count := services.Tag.Paginate(r.IndexValidate(c))

	responses.Json(c, tags, count)
}

func (*tagController) Store(c *gin.Context) {
	responses.Json(c, r.StoreValidate(c).Store())
}

func (*tagController) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id <= 0 {
		responses.Json(c, nil, false, http.StatusNotFound, responses.CodeNotFound)
		return
	}

	responses.Json(c, r.UpdateValidate(c).Update(id))
}
