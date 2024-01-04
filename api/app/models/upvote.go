package models

import (
	"blog/app"
	"gorm.io/gorm"
	"time"
)

type Upvote struct {
	ID        int64           `gorm:"primaryKey;type:bigint unsigned" json:"id"`
	UserID    int64           `gorm:"index:idx_unique_user_upvote,unique,priority:2;index;type:bigint unsigned;not null" json:"userId"`
	CommentID int64           `gorm:"index:idx_unique_user_upvote,unique,priority:1;type:bigint unsigned;not null" json:"commentId"`
	CreatedAt *time.Time      `gorm:"not null" json:"createdAt,omitempty"`
	UpdatedAt *time.Time      `gorm:"not null" json:"updatedAt,omitempty"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	User      *User           `json:"user,omitempty"`
	Comment   *Comment        `json:"comment,omitempty"`
}

type Upvotes []*Upvote

func (upvote *Upvote) Load(preloads map[string][]any, missing ...bool) {
	LoadRelations(app.DB, upvote, preloads, missing...)
}

func (upvotes Upvotes) Load(preloads map[string][]any, missing ...bool) {
	LoadRelations(app.DB, upvotes, preloads, missing...)
}

func (upvotes Upvotes) ParseInclude(relations any) {
	availableRelations := map[string][]any{"Comment.User": nil, "Comment.Article": nil}

	upvotes.Load(FilterRelations(relations, availableRelations))
}
