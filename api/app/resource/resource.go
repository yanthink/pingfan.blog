package resource

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

type Type string

const (
	Avatar       Type = "avatar"
	ArticleImage Type = "articleImage"
	CommentImage Type = "commentImage"
)

func (t *Type) Validate(file *multipart.FileHeader) {
	switch *t {
	case Avatar:
		NewDimensions().
			Accept("image/*").
			MinWidth(200).
			Ratio(1).
			MaxSize(1 << 20).
			Validate(file)
	case ArticleImage:
		NewDimensions().
			Accept("image/*").
			MinWidth(200).
			MaxSize(2 << 20).
			Validate(file)
	case CommentImage:
		NewDimensions().
			Accept("image/*").
			MinWidth(50).
			MaxSize(512 << 10).
			Validate(file)
	}
}

func (t *Type) UploadPath(date ...time.Time) string {
	d := time.Now()

	if len(date) > 0 {
		d = date[0]
	}

	return fmt.Sprintf("upload/%s", d.Format("20060102"))
}

func (t *Type) StorePath() string {
	switch *t {
	case Avatar:
		return fmt.Sprintf("avatars/%s", time.Now().Format("20060102"))
	case ArticleImage:
		return fmt.Sprintf("articles/%s", time.Now().Format("20060102"))
	case CommentImage:
		return fmt.Sprintf("comments/%s", time.Now().Format("20060102"))
	}

	return ""
}

func (t *Type) IsUploadPath(path string) bool {
	return strings.HasPrefix(path, filepath.Dir(t.UploadPath()))
}

func (t *Type) Resize() (resize [2]int) {
	switch *t {
	case Avatar:
		resize = [2]int{200, 200}
	}

	return
}
