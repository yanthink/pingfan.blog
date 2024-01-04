package notifications

import (
	"blog/app/models"
	"fmt"
)

// CommentHasNewReply 文章有新回复通知
type CommentHasNewReply struct {
	notification
	Comment *models.Comment
	Article *models.Article
}

func (c *CommentHasNewReply) Setup() error {
	if c.Comment.ParentID != c.Comment.CommentID {
		c.Comment.Article = c.Article
		c.Comment.Load(map[string][]any{"User": nil, "Article": nil, "Comment.User": nil, "Parent": nil})
	}

	return nil
}

func (c *CommentHasNewReply) Users() (users models.Users) {
	if c.Comment.ParentID != c.Comment.CommentID && c.Comment.UserID != c.Comment.Comment.UserID && c.Comment.Parent.UserID != c.Comment.Comment.UserID {
		users = models.Users{c.Comment.Comment.User}
	}

	return
}

func (c *CommentHasNewReply) FromUser() (user *models.User) {
	return c.Comment.User
}

func (c *CommentHasNewReply) Subject() string {
	articleUrl := c.Comment.Article.Url()
	articleUrl.Fragment = "评论"

	query := articleUrl.Query()
	query.Add("pinnedId", fmt.Sprintf("%d", c.Comment.ID))
	articleUrl.RawQuery = query.Encode()

	return fmt.Sprintf("你的评论：[%s](%s) 有新回复", c.Comment.Article.Title, articleUrl)
}

func (c *CommentHasNewReply) Message() string {
	return c.Comment.Content
}

func (c *CommentHasNewReply) ToDatabase() (*map[string]any, bool) {
	data := map[string]any{
		"id":         c.Comment.ID,
		"article_id": c.Comment.ArticleID,
		"comment_id": c.Comment.CommentID,
	}

	return &data, true
}
