package filters

import "gorm.io/gorm"

type ArticleTagFilter struct {
}

func (f *ArticleTagFilter) ID(db *gorm.DB, id int64, _ any) *gorm.DB {
	return db.Where("tag_id = ?", id)
}
