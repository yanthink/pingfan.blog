package models

import "time"

type Tag struct {
	ID        int64      `gorm:"primaryKey;type:bigint unsigned" form:"id" json:"id"`
	Name      string     `gorm:"size:40;not null" json:"name"`
	Sort      int64      `gorm:"type:smallint unsigned;not null;default:0" json:"sort,omitempty"`
	CreatedAt *time.Time `gorm:"not null" json:"createdAt,omitempty"`
	UpdatedAt *time.Time `gorm:"not null" json:"updatedAt,omitempty"`
}

type Tags []*Tag
