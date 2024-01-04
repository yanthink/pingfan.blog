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

// CommentUpvote 评论点赞通知
type CommentUpvote struct {
	notification
	Upvote  *models.Upvote
	Comment *models.Comment
}

func (c *CommentUpvote) Setup() error {
	if c.Upvote.ID > 0 && c.Upvote.UpdatedAt != nil && c.Upvote.CreatedAt.Equal(*c.Upvote.UpdatedAt) {
		c.Upvote.Comment = c.Comment
		c.Upvote.Load(map[string][]any{"User": nil, "Comment.User": nil, "Comment.Article": nil})
	}

	return nil
}

func (c *CommentUpvote) Users() (users models.Users) {
	if c.Upvote.ID > 0 && c.Upvote.UpdatedAt != nil && c.Upvote.CreatedAt.Equal(*c.Upvote.UpdatedAt) && c.Upvote.UserID != c.Upvote.Comment.UserID {
		users = models.Users{c.Upvote.Comment.User}
	}

	return
}

func (c *CommentUpvote) FromUser() (user *models.User) {
	return c.Upvote.User
}

func (c *CommentUpvote) Subject() string {
	articleUrl := c.Upvote.Comment.Article.Url()
	articleUrl.Fragment = "评论"

	query := articleUrl.Query()
	query.Add("pinnedId", fmt.Sprintf("%d", c.Upvote.CommentID))
	articleUrl.RawQuery = query.Encode()

	return fmt.Sprintf(
		"[%s](%s) • 赞了您的评论：[%s](%s)",
		c.Upvote.User.Name,
		c.Upvote.User.Url(),
		c.Upvote.Comment.Article.Title,
		articleUrl,
	)
}

func (c *CommentUpvote) Message() string {
	var buf bytes.Buffer
	_ = goldmark.Convert([]byte(c.Comment.Content), &buf)

	pc := regexp.MustCompile(`<img[^>]*>`).ReplaceAllString(buf.String(), "[图片]")
	pc = regexp.MustCompile(`\r?\n`).ReplaceAllString(pc, " ")

	qc := []rune(bluemonday.StrictPolicy().Sanitize(pc))
	if len(qc) > 200 {
		qc = qc[:200]
	}

	commentUrl := c.Comment.Article.Url()
	commentUrl.Fragment = "评论"
	query := commentUrl.Query()
	query.Add("cid", fmt.Sprintf("%d", c.Comment.ID))
	query.Add("parentId", fmt.Sprintf("%d", c.Comment.ParentID))
	commentUrl.RawQuery = query.Encode()

	if qc := strings.TrimSpace(string(qc)); len(qc) > 0 {
		return fmt.Sprintf(`<p class="reply-quote"><a href="%s">%s</a></p>`, commentUrl, qc)
	}

	return ""
}

func (c *CommentUpvote) ToDatabase() (*map[string]any, bool) {
	data := map[string]any{
		"id":         c.Upvote.ID,
		"comment_id": c.Upvote.CommentID,
	}

	return &data, true
}
