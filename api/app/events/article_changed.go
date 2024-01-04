package events

import "blog/app/models"

type ArticleChanged struct {
	Original *models.Article
	Article  *models.Article
}
