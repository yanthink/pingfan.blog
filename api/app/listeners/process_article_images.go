package listeners

import (
	"blog/app"
	"blog/app/events"
	"blog/app/helpers"
	"blog/app/models"
	"blog/app/resource"
	"blog/app/services"
	"blog/config"
	"fmt"
	"strings"
)

type ProcessArticleImages struct {
	listener
}

func (*ProcessArticleImages) Handle(event any) (err error) {
	e, ok := event.(*events.ArticleChanged)
	if !ok {
		return nil
	}

	original := e.Original
	article := e.Article

	if article == nil {
		return
	}

	if original == nil {
		original = &models.Article{}
	}

	if article.Preview == original.Preview && article.Content == original.Content {
		return
	}

	lockKey := fmt.Sprintf("%s_article_update_lock:%d", config.Redis.Prefix, article.ID)
	mutex := app.Redsync.NewMutex(lockKey)
	if mutex.Lock() != nil {
		return nil
	}
	defer mutex.Unlock()

	preview := services.Resource.Sync([]string{article.Preview}, []string{original.Preview}, resource.ArticleImage)[0]

	content := article.Content

	newUrls := helpers.MatchImageUrls(content)
	oldUrls := helpers.MatchImageUrls(original.Content)

	urls := services.Resource.Sync(newUrls, oldUrls, resource.ArticleImage)

	for i, url := range urls {
		if url != newUrls[i] {
			content = strings.ReplaceAll(content, newUrls[i], url)
		}
	}

	if preview == article.Preview && content == article.Content {
		return
	}

	article.Preview = preview
	article.Content = content

	app.DB.Model(&article).Select("Preview", "Content").UpdateColumns(article)

	return
}
