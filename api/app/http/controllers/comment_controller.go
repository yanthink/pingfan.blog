package controllers

import (
	r "blog/app/http/requests/comments"
	"blog/app/http/responses"
	"blog/app/models"
	"blog/app/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type commentController struct {
}

func (*commentController) Index(c *gin.Context) {
	req := r.IndexValidate(c)
	comments, count := services.Comment.Paginate(req)

	responses.Json(c, &comments, count, func(*responses.JsonResponse) {
		flatComments := comments.Flat()

		if req.WrapID > 0 {
			comment := services.Comment.GetByID(req.WrapID)
			flatComments = append(flatComments, comment)

			comment.Load(map[string][]any{"User": nil})
			comment.Replies = comments

			comments = models.Comments{comment}
		}

		if req.WithReplyUser {
			flatComments.Load(map[string][]any{
				"Parent": {func(db *gorm.DB) *gorm.DB {
					return db.Where("comment_id > 0")
				}},
				"Parent.User": nil,
			})
		}

		if req.PinnedID > 0 && req.Page == 1 {
			comment := services.Comment.GetByID(req.PinnedID)
			comment.Load(map[string][]any{"User": nil})
			comments = append(models.Comments{comment}, comments...)
		}

		flatComments.WithUserHasUpvoted(c.GetInt64("userId"))
	})
}

func (*commentController) CursorPaginate(c *gin.Context) {
	comments, cursor := services.Comment.CursorPaginate(r.CursorPaginateValidate(c))
	comments.WithUserHasUpvoted(c.GetInt64("userId"))

	responses.Json(c, comments, cursor)
}

func (*commentController) Show(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id <= 0 {
		responses.Json(c, nil, false, http.StatusNotFound, responses.CodeNotFound)
		return
	}

	responses.Json(c, services.Comment.GetByID(id))
}

func (*commentController) Store(c *gin.Context) {
	responses.Json(c, r.StoreValidate(c).Store())
}

func (*commentController) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id <= 0 {
		responses.Json(c, nil, false, http.StatusNotFound, responses.CodeNotFound)
		return
	}

	responses.Json(c, r.UpdateValidate(c).Update(id))
}

func (*commentController) Upvote(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id <= 0 {
		responses.Json(c, nil, false, http.StatusNotFound, responses.CodeNotFound)
		return
	}

	userId := c.GetInt64("userId")
	comment := services.Comment.Upvote(id, userId)

	responses.Json(c, &models.Comment{
		ID:          comment.ID,
		UpvoteCount: comment.UpvoteCount,
		HasUpvoted:  comment.HasUpvoted,
	})
}
