package filters

import (
	"gorm.io/gorm"
	"strings"
)

type CommentFilter struct {
}

func (f *CommentFilter) ID(db *gorm.DB, id int64) *gorm.DB {
	return db.Where("id = ?", id)
}

func (f *CommentFilter) UserID(db *gorm.DB, id int64) *gorm.DB {
	return db.Where("user_id = ?", id)
}

func (f *CommentFilter) ArticleID(db *gorm.DB, id int64) *gorm.DB {
	return db.Where("article_id = ?", id)
}

func (f *CommentFilter) CommentID(db *gorm.DB, id *int64) *gorm.DB {
	return db.Where("comment_id = ?", id)
}

func (f *CommentFilter) ParentID(db *gorm.DB, id *int64) *gorm.DB {
	return db.Where("parent_id = ?", id)
}

func (f *CommentFilter) Sort(db *gorm.DB, sort string) *gorm.DB {
	if strings.ToUpper(sort) == "DESC" {
		return db.Order("id DESC")
	}

	return db
}

func (f *CommentFilter) Type(db *gorm.DB, t *int64) *gorm.DB {
	if *t == 0 {
		return db.Where("comment_id = ?", 0)
	}

	return db.Where("comment_id > ?", 0)
}
