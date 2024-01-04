package services

import (
	"blog/app"
	"blog/app/filters"
	"blog/app/models"
	"blog/app/pagination"
)

type favoriteService struct {
}

func (*favoriteService) Paginate(paginator pagination.Pager) (favorites models.Favorites, count int64) {
	tx := app.DB.
		Model(favorites).
		Order("updated_at DESC").
		Order("id DESC").
		Scopes(filters.New(&filters.FavoriteFilter{}, paginator))

	var err error

	if _, count, err = paginator.Paginate(tx, &favorites); err != nil {
		panic(err)
	}

	return
}
