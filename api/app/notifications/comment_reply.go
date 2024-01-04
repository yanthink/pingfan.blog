package notifications

import (
	"blog/app/models"
	"bytes"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"regexp"
	"strings"
)

type CommentReply struct {
	notification
	Comment *models.Comment
	Article *models.Article
}

func (c *CommentReply) Setup() error {
	if c.Comment.ParentID > 0 {
		c.Comment.Article = c.Article
		c.Comment.Load(map[string][]any{"User": nil, "Article": nil, "Parent.User": nil})
	}

	return nil
}

func (c *CommentReply) Users() (users models.Users) {
	if c.Comment.ParentID > 0 && c.Comment.UserID != c.Comment.Parent.UserID {
		users = models.Users{c.Comment.Parent.User}
	}

	return
}

func (c *CommentReply) FromUser() (user *models.User) {
	return c.Comment.User
}

func (c *CommentReply) Subject() string {
	articleUrl := c.Comment.Article.Url()
	articleUrl.Fragment = "评论"

	query := articleUrl.Query()
	query.Add("parentId", fmt.Sprintf("%d", c.Comment.ID))
	query.Add("cid", fmt.Sprintf("%d", c.Comment.ID))
	query.Add("parentId", fmt.Sprintf("%d", c.Comment.ParentID))
	articleUrl.RawQuery = query.Encode()

	return fmt.Sprintf(
		"[%s](%s) • 回复了您的评论：[%s](%s)",
		c.Comment.User.Name,
		c.Comment.User.Url(),
		c.Comment.Article.Title,
		articleUrl,
	)
}

func (c *CommentReply) Message() string {
	var buf bytes.Buffer
	_ = goldmark.Convert([]byte(c.Comment.Parent.Content), &buf)

	pc := regexp.MustCompile(`<img[^>]*>`).ReplaceAllString(buf.String(), "[图片]")
	pc = regexp.MustCompile(`\r?\n`).ReplaceAllString(pc, " ")

	qc := []rune(bluemonday.StrictPolicy().Sanitize(pc))
	if len(qc) > 200 {
		qc = qc[:200]
	}

	replyUrl := c.Comment.Article.Url()
	replyUrl.Fragment = "评论"
	query := replyUrl.Query()
	query.Add("cid", fmt.Sprintf("%d", c.Comment.ID))
	query.Add("parentId", fmt.Sprintf("%d", c.Comment.ParentID))
	replyUrl.RawQuery = query.Encode()

	if qc := strings.TrimSpace(string(qc)); len(qc) > 0 {
		return fmt.Sprintf(`%s\n\n<p class="reply-quote"><a href="%s">%s</a></p>`, c.Comment.Content, replyUrl, qc)
	}

	return c.Comment.Content
}

func (c *CommentReply) ToDatabase() (*map[string]any, bool) {
	data := map[string]any{
		"id":         c.Comment.ID,
		"article_id": c.Comment.ArticleID,
		"comment_id": c.Comment.CommentID,
	}

	return &data, true
}
