package models

import (
	"blog/app"
	"gorm.io/gorm"
	"time"
)

type Favorite struct {
	ID        int64           `gorm:"primaryKey;type:bigint unsigned" json:"id"`
	UserID    int64           `gorm:"uniqueIndex:user_favorite_unique_idx,priority:2;type:bigint unsigned;not null" json:"userId"`
	ArticleID int64           `gorm:"uniqueIndex:user_favorite_unique_idx,priority:1;type:bigint unsigned;not null" json:"articleId"`
	CreatedAt *time.Time      `gorm:"not null" json:"createdAt,omitempty"`
	UpdatedAt *time.Time      `gorm:"not null" json:"updatedAt,omitempty"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	Article   *Article        `json:"article,omitempty"`
}

type Favorites []*Favorite

func (favorite *Favorite) Load(preloads map[string][]any, missing ...bool) {
	LoadRelations(app.DB, favorite, preloads, missing...)
}

func (favorites Favorites) Load(preloads map[string][]any, missing ...bool) {
	LoadRelations(app.DB, favorites, preloads, missing...)
}

func (favorites Favorites) ParseInclude(relations any) {
	availableRelations := map[string][]any{"Article.Tags": nil, "Article.User": nil}

	favorites.Load(FilterRelations(relations, availableRelations))
}
