package models

import (
	"blog/config"
	"fmt"
	"gorm.io/gorm"
	"net/url"
	"time"
)

type UserRole byte

const (
	UserRoleNormal UserRole = iota
	UserRoleManage
)

type User struct {
	ID          int64           `gorm:"primaryKey;type:bigint unsigned" json:"id"`
	Name        *string         `gorm:"unique;size:20" json:"name"`
	Email       string          `gorm:"unique;size:255" json:"email,omitempty"`
	Openid      *string         `gorm:"unique;size:40" json:"-"`
	Password    string          `gorm:"size:60;not null;default:''" json:"-"`
	Avatar      string          `gorm:"size:255;not null;default:''" json:"avatar,omitempty"`
	Role        UserRole        `gorm:"type:tinyint unsigned;not null;default:0;comment:角色- 0：普通用户，1：管理用户" json:"role,omitempty"`
	Status      *int64          `gorm:"type:tinyint unsigned;not null;default:0;comment:状态- 0：正常，1：锁定" json:"status,omitempty"`
	Meta        *map[string]any `gorm:"type:json;serializer:json" json:"meta,omitempty"`
	HasPassword bool            `gorm:"-" json:"hasPassword,omitempty"`
	CreatedAt   *time.Time      `gorm:"index;not null" json:"createdAt,omitempty"`
	UpdatedAt   *time.Time      `gorm:"not null" json:"updatedAt,omitempty"`
}

type Users []*User

func (user *User) AfterFind(_ *gorm.DB) (err error) {
	user.HasPassword = user.Password != ""

	return
}

func (user *User) Url() *url.URL {
	parsed, _ := url.Parse(fmt.Sprintf("%s/users/%d", config.App.SiteUrl, user.ID))

	return parsed
}
