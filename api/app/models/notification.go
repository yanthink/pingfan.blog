package models

import (
	"blog/app"
	"time"
)

type Notification struct {
	ID         int64           `gorm:"primaryKey;type:bigint unsigned" json:"id"`
	UserID     int64           `gorm:"index:idx_user_read_at,priority:1;type:bigint unsigned;not null" json:"userId"`
	FromUserID int64           `gorm:"type:bigint unsigned;not null" json:"fromUserId"`
	Type       string          `gorm:"size:100;not null" json:"type"`
	Subject    string          `gorm:"size:255;not null" json:"subject"`
	Message    string          `gorm:"size:2048;not null" json:"message"`
	Data       *map[string]any `gorm:"type:json;serializer:json" json:"data,omitempty"`
	ReadAt     *time.Time      `gorm:"index:idx_user_read_at,priority:2;null" json:"readAt,omitempty"`
	CreatedAt  *time.Time      `gorm:"not null" json:"createdAt,omitempty"`
	UpdatedAt  *time.Time      `gorm:"not null" json:"updatedAt,omitempty"`
	User       *User           `json:"user,omitempty"`
	FromUser   *User           `json:"fromUser,omitempty"`
}

type Notifications []*Notification

func (notification *Notification) Load(preloads map[string][]any, missing ...bool) {
	LoadRelations(app.DB, notification, preloads, missing...)
}

func (notifications Notifications) Load(preloads map[string][]any, missing ...bool) {
	LoadRelations(app.DB, notifications, preloads, missing...)
}

func (notifications Notifications) ParseInclude(relations any) {
	availableRelations := map[string][]any{"FromUser": nil}

	notifications.Load(FilterRelations(relations, availableRelations))
}
