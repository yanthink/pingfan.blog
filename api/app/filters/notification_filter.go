package filters

import (
	"gorm.io/gorm"
)

type NotificationFilter struct {
}

func (f *NotificationFilter) UserID(db *gorm.DB, id int64, _ any) *gorm.DB {
	return db.Where("user_id = ?", id)
}

func (f *NotificationFilter) Type(db *gorm.DB, t string, _ any) *gorm.DB {
	return db.Where("type = ?", t)
}
