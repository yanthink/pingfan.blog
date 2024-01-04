package notifications

import (
	"blog/app/models"
	"fmt"
)

// ArticleLike 文章点赞通知
type ArticleLike struct {
	notification
	Like    *models.Like
	Article *models.Article
}

func (c *ArticleLike) Setup() error {
	if c.Like.ID > 0 && c.Like.UpdatedAt != nil && c.Like.CreatedAt.Equal(*c.Like.UpdatedAt) {
		c.Like.Article = c.Article
		c.Like.Load(map[string][]any{"User": nil, "Article.User": nil})
	}

	return nil
}

func (c *ArticleLike) Users() (users models.Users) {
	if c.Like.ID > 0 && c.Like.UpdatedAt != nil && c.Like.CreatedAt.Equal(*c.Like.UpdatedAt) && c.Like.UserID != c.Like.Article.UserID {
		users = models.Users{c.Like.Article.User}
	}

	return
}

func (c *ArticleLike) FromUser() (user *models.User) {
	return c.Like.User
}

func (c *ArticleLike) Subject() string {
	articleUrl := c.Like.Article.Url()

	return fmt.Sprintf(
		"[%s](%s) • 赞了你的文章：[%s](%s)",
		c.Like.User.Name,
		c.Like.User.Url(),
		c.Like.Article.Title,
		articleUrl,
	)
}

func (c *ArticleLike) ToDatabase() (*map[string]any, bool) {
	data := map[string]any{
		"id":         c.Like.ID,
		"article_id": c.Like.ArticleID,
	}

	return &data, true
}
