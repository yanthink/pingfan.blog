package filters

import (
	"gorm.io/gorm"
	"strings"
)

type LikeFilter struct {
}

func (f *LikeFilter) UserID(db *gorm.DB, id int64) *gorm.DB {
	return db.Where("user_id = ?", id)
}

func (f *LikeFilter) ArticleID(db *gorm.DB, id int64) *gorm.DB {
	return db.Where("article_id = ?", id)
}

func (f *LikeFilter) Sort(db *gorm.DB, sort string) *gorm.DB {
	if strings.ToUpper(sort) == "DESC" {
		return db.Order("id DESC")
	}

	return db
}
