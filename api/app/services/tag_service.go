package services

import (
	"blog/app"
	"blog/app/filters"
	"blog/app/models"
	"blog/app/pagination"
	"gorm.io/gorm/clause"
)

type tagService struct {
}

func (*tagService) Paginate(paginator pagination.Pager) (tags models.Tags, count int64) {
	tx := app.DB.
		Model(tags).
		Order("sort ASC").
		Scopes(filters.New(&filters.TagFilter{}, paginator))

	var err error

	if _, count, err = paginator.Paginate(tx, &tags); err != nil {
		panic(err)
	}

	return
}

func (*tagService) Add(tag *models.Tag) *models.Tag {
	if err := app.DB.Create(&tag).Error; err != nil {
		panic(err)
	}

	return tag
}

func (*tagService) Update(id int64, tag *models.Tag) *models.Tag {
	tag.ID = id

	if err := app.DB.Model(tag).Omit(clause.Associations).Updates(&tag).Error; err != nil {
		panic(err)
	}

	return tag
}
