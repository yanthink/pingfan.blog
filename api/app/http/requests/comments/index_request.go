package comments

import (
	"blog/app/pagination"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
)

type IndexRequest struct {
	pagination.Paginator
	ID            int64  `form:"cid" json:"cid"`
	ArticleID     int64  `form:"articleId" json:"articleId" binding:"omitempty,gt=0"`
	CommentID     *int64 `form:"commentId" json:"commentId" binding:"omitempty,gte=0"`
	ParentID      *int64 `form:"parentId" json:"parentId" binding:"omitempty,gte=0"`
	PinnedID      int64  `form:"pinnedId" json:"pinnedId" label:"置顶ID"`
	WrapID        int64  `form:"wrapId" json:"wrapId" label:"包裹ID"`
	WithReplyUser bool   `form:"withReplyUser" json:"withReplyUser"`
	Sort          string `form:"sort" json:"sort"`
}

func IndexValidate(c *gin.Context) (r *IndexRequest) {
	r = &IndexRequest{}
	validation.Validate(c, r)

	return
}
