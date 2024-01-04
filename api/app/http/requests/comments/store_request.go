package comments

import (
	h "blog/app/http"
	"blog/app/http/responses"
	"blog/app/mini_program"
	"blog/app/models"
	"blog/app/services"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

type StoreRequest struct {
	ArticleID int64  `form:"articleId" json:"articleId" binding:"required,gt=0"`
	Content   string `form:"title" json:"content" label:"内容" binding:"required,max=512"`
	ParentID  int64  `form:"parentId" json:"parentId"`
	CommentID int64  `form:"-" json:"-"`
	UserID    int64  `form:"-" json:"-"`
}

func (r *StoreRequest) Store() *models.Comment {
	comment := models.Comment{
		ArticleID: r.ArticleID,
		UserID:    r.UserID,
		Content:   r.Content,
		ParentID:  r.ParentID,
		CommentID: r.CommentID,
	}

	return services.Comment.Add(&comment)
}

func StoreValidate(c *gin.Context) (r *StoreRequest) {
	r = &StoreRequest{}
	validation.Validate(c, r)

	if r.ParentID > 0 {
		replyComment := services.Comment.GetByID(r.ParentID)

		if replyComment.ArticleID != r.ArticleID {
			panic(&h.Error{StatusCode: http.StatusForbidden, Code: responses.CodeAccessDenied, Message: "非法操作"})
		}

		if r.CommentID = replyComment.CommentID; r.CommentID == 0 {
			r.CommentID = replyComment.ID
		}
	}

	if user := services.User.GetAuthUser(c); user.Openid != nil {
		if !mini_program.Wx.MsgSecCheck(r.Content, *user.Openid) {
			panic(&validation.Error{
				Message: "含有违法违规内容！",
				Errors: map[string][]string{
					"content": {"含有违法违规内容！"},
				},
			})
		}
	}

	r.UserID = c.GetInt64("userId")

	return
}
