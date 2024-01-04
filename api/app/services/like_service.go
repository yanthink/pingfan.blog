package services

import (
	"blog/app"
	"blog/app/filters"
	"blog/app/models"
	"blog/app/pagination"
)

type likeService struct {
}

func (*likeService) Paginate(paginator pagination.Pager) (likes models.Likes, count int64) {
	tx := app.DB.
		Model(likes).
		Scopes(filters.New(&filters.LikeFilter{}, paginator))

	var err error

	if _, count, err = paginator.Paginate(tx, &likes); err != nil {
		panic(err)
	}

	return
}
