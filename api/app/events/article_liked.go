package events

import "blog/app/models"

type ArticleLiked struct {
	Like    *models.Like
	Article *models.Article
}
