package listeners

import (
	"blog/app"
	"blog/app/events"
	"blog/app/helpers"
	"blog/app/resource"
	"blog/app/services"
	"strings"
)

type ProcessCommentImages struct {
	listener
}

func (*ProcessCommentImages) Handle(event any) (err error) {
	e, ok := event.(*events.Commented)
	if !ok {
		return nil
	}

	comment := e.Comment

	content := comment.Content

	newUrls := helpers.MatchImageUrls(content)
	var oldUrls []string

	urls := services.Resource.Sync(newUrls, oldUrls, resource.ArticleImage)

	for i, url := range urls {
		if url != newUrls[i] {
			content = strings.ReplaceAll(content, newUrls[i], url)
		}
	}

	if content == comment.Content {
		return
	}

	comment.Content = content

	app.DB.Model(&comment).Select("Content").UpdateColumns(comment)

	return
}
