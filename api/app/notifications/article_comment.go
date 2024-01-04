package notifications

import (
	"blog/app/models"
	"fmt"
)

// ArticleComment 文章评论通知
type ArticleComment struct {
	notification
	Comment *models.Comment
	Article *models.Article
}

func (c *ArticleComment) Setup() error {
	if c.Comment.ParentID == 0 {
		c.Comment.Article = c.Article
		c.Comment.Load(map[string][]any{"User": nil, "Article.User": nil})
	}

	return nil
}

func (c *ArticleComment) Users() (users models.Users) {
	if c.Comment.ParentID == 0 && c.Comment.UserID != c.Comment.Article.UserID {
		users = models.Users{c.Comment.Article.User}
	}

	return
}

func (c *ArticleComment) FromUser() (user *models.User) {
	return c.Comment.User
}

func (c *ArticleComment) Subject() string {
	articleUrl := c.Comment.Article.Url()
	articleUrl.Fragment = "评论"

	query := articleUrl.Query()
	query.Add("pinnedId", fmt.Sprintf("%d", c.Comment.ID))
	articleUrl.RawQuery = query.Encode()

	return fmt.Sprintf(
		"[%s](%s) • 评论了你的文章：[%s](%s)",
		c.Comment.User.Name,
		c.Comment.User.Url(),
		c.Comment.Article.Title,
		articleUrl,
	)
}

func (c *ArticleComment) Message() string {
	return c.Comment.Content
}

func (c *ArticleComment) ToDatabase() (*map[string]any, bool) {
	data := map[string]any{
		"id":         c.Comment.ID, // 评论id
		"article_id": c.Comment.ArticleID,
	}

	return &data, true
}
