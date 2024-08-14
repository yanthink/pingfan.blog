package filters

import (
	"gorm.io/gorm"
	"strings"
)

type UpvoteFilter struct {
}

func (f *UpvoteFilter) UserID(db *gorm.DB, id int64, _ any) *gorm.DB {
	return db.Where("user_id = ?", id)
}

func (f *UpvoteFilter) CommentID(db *gorm.DB, id int64, _ any) *gorm.DB {
	return db.Where("comment_id = ?", id)
}

func (f *UpvoteFilter) Sort(db *gorm.DB, sort string, _ any) *gorm.DB {
	if strings.ToUpper(sort) == "DESC" {
		return db.Order("id DESC")
	}

	return db
}
