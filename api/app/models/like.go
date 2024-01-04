package models

import (
	"blog/app"
	"gorm.io/gorm"
	"time"
)

type Like struct {
	ID        int64           `gorm:"primaryKey;type:bigint unsigned" json:"id"`
	UserID    int64           `gorm:"index:idx_unique_user_like,unique,priority:2;index;type:bigint unsigned;not null" json:"userId"`
	ArticleID int64           `gorm:"index:idx_unique_user_like,unique,priority:1;type:bigint unsigned;not null" json:"ArticleID"`
	CreatedAt *time.Time      `gorm:"not null" json:"createdAt,omitempty"`
	UpdatedAt *time.Time      `gorm:"not null" json:"updatedAt,omitempty"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	User      *User           `json:"user,omitempty"`
	Article   *Article        `json:"article,omitempty"`
}

type Likes []*Like

func (like *Like) Load(preloads map[string][]any, missing ...bool) {
	LoadRelations(app.DB, like, preloads, missing...)
}

func (likes Likes) Load(preloads map[string][]any, missing ...bool) {
	LoadRelations(app.DB, likes, preloads, missing...)
}

func (likes Likes) ParseInclude(relations any) {
	availableRelations := map[string][]any{"Article.Tags": nil, "Article.User": nil}

	likes.Load(FilterRelations(relations, availableRelations))
}
