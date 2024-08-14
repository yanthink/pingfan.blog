package filters

import (
	"blog/app/models"
	"gorm.io/gorm"
)

type UserFilter struct {
}

func (f *UserFilter) ID(db *gorm.DB, id int64, _ any) *gorm.DB {
	return db.Where("id = ?", id)
}

func (f *UserFilter) Name(db *gorm.DB, name string, _ any) *gorm.DB {
	return db.Where("name = ?", name)
}

func (f *UserFilter) Email(db *gorm.DB, email string, _ any) *gorm.DB {
	return db.Where("email = ?", email)
}

func (f *UserFilter) Openid(db *gorm.DB, openid string, _ any) *gorm.DB {
	return db.Where("openid = ?", openid)
}

func (f *UserFilter) Role(db *gorm.DB, role *models.UserRole, _ any) *gorm.DB {
	return db.Where("role = ?", role)
}

func (f *UserFilter) Status(db *gorm.DB, status *int64, _ any) *gorm.DB {
	return db.Where("status = ?", status)
}
