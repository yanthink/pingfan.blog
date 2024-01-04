package services

import (
	"blog/app"
	"blog/app/filters"
	"blog/app/models"
	"blog/app/pagination"
)

type upvoteService struct {
}

func (*upvoteService) Paginate(paginator pagination.Pager) (upvotes models.Upvotes, count int64) {
	tx := app.DB.
		Model(upvotes).
		Scopes(filters.New(&filters.UpvoteFilter{}, paginator))

	var err error

	if _, count, err = paginator.Paginate(tx, &upvotes); err != nil {
		panic(err)
	}

	return
}
