package controllers

import (
	r "blog/app/http/requests/articles"
	"blog/app/http/responses"
	"blog/app/models"
	"blog/app/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type articleController struct {
}

func (*articleController) Index(c *gin.Context) {
	articles, count := services.Article.Paginate(r.IndexValidate(c))
	responses.Json(c, articles, count)
}

func (*articleController) Search(c *gin.Context) {
	req := r.SearchValidate(c)

	articles, count := services.Article.SearchPaginate(strings.TrimSpace(req.Q), req.Filters, req.QueryFields, req.SortFields, req.Limit, req.Page)
	responses.Json(c, articles, count)
}

func (*articleController) CursorPaginate(c *gin.Context) {
	articles, cursor := services.Article.CursorPaginate(r.CursorPaginateValidate(c))

	responses.Json(c, articles, cursor)
}

func (*articleController) Show(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id <= 0 {
		responses.Json(c, nil, false, http.StatusNotFound, responses.CodeNotFound)

		return
	}

	userId := c.GetInt64("userId")

	article := services.Article.GetByID(id)
	article.WithUserHasLiked(userId)
	article.WithUserHasFavorited(userId)
	article.IncrViewCount(1)

	responses.Json(c, article)
}

func (*articleController) Store(c *gin.Context) {
	responses.Json(c, r.StoreValidate(c).Store())
}

func (*articleController) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id <= 0 {
		responses.Json(c, nil, false, http.StatusNotFound, responses.CodeNotFound)
		return
	}

	responses.Json(c, r.UpdateValidate(c, id).Update(id))
}

func (*articleController) Like(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id <= 0 {
		responses.Json(c, nil, false, http.StatusNotFound, responses.CodeNotFound)
		return
	}

	userId := c.GetInt64("userId")
	article := services.Article.Like(id, userId)

	responses.Json(c, &models.Article{
		ID:        article.ID,
		LikeCount: article.LikeCount,
		HasLiked:  article.HasLiked,
	})
}

func (*articleController) Favorite(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id <= 0 {
		responses.Json(c, nil, false, http.StatusNotFound, responses.CodeNotFound)
		return
	}

	userId := c.GetInt64("userId")
	article := services.Article.Favorite(id, userId)

	responses.Json(c, &models.Article{
		ID:            article.ID,
		FavoriteCount: article.FavoriteCount,
		HasFavorited:  article.HasFavorited,
	})
}
