package filters

import (
	"gorm.io/gorm"
)

type FavoriteFilter struct {
}

func (f *FavoriteFilter) UserID(db *gorm.DB, id int64, _ any) *gorm.DB {
	return db.Where("user_id = ?", id)
}
