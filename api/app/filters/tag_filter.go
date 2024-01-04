package filters

import (
	"fmt"
	"gorm.io/gorm"
)

type TagFilter struct {
}

func (f *TagFilter) Q(db *gorm.DB, q string) *gorm.DB {
	return db.Where("name like ?", fmt.Sprintf("%%%s%%", q))
}
