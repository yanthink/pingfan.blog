package notifications

import (
	"blog/app/models"
	"fmt"
)

// ArticleHasNewReply 文章有新回复通知
type ArticleHasNewReply struct {
	notification
	Comment *models.Comment
	Article *models.Article
}

func (c *ArticleHasNewReply) Setup() error {
	if c.Comment.ParentID > 0 {
		c.Comment.Article = c.Article
		c.Comment.Load(map[string][]any{"User": nil, "Article.User": nil, "Parent": nil, "Comment": nil})
	}

	return nil
}

func (c *ArticleHasNewReply) Users() (users models.Users) {
	if c.Comment.ParentID > 0 &&
		c.Comment.UserID != c.Comment.Article.UserID &&
		c.Comment.Parent.UserID != c.Comment.Article.UserID &&
		c.Comment.Comment.UserID != c.Comment.Article.UserID {
		users = models.Users{c.Comment.Article.User}
	}

	return
}

func (c *ArticleHasNewReply) FromUser() (user *models.User) {
	return c.Comment.User
}

func (c *ArticleHasNewReply) Subject() string {
	articleUrl := c.Comment.Article.Url()
	articleUrl.Fragment = "评论"

	query := articleUrl.Query()
	query.Add("pinnedId", fmt.Sprintf("%d", c.Comment.ID))
	articleUrl.RawQuery = query.Encode()

	return fmt.Sprintf("你的文章：[%s](%s) 有新回复", c.Comment.Article.Title, articleUrl)
}

func (c *ArticleHasNewReply) Message() string {
	return c.Comment.Content
}

func (c *ArticleHasNewReply) ToDatabase() (*map[string]any, bool) {
	data := map[string]any{
		"id":         c.Comment.ID,
		"article_id": c.Comment.ArticleID,
		"comment_id": c.Comment.CommentID,
	}

	return &data, true
}
