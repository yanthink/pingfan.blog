package models

import (
	"blog/app"
	"blog/config"
	"bytes"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"gorm.io/gorm"
	"math"
	"net/url"
	"time"
)

const (
	ArticleViewHotnessWeight        = 1        // 浏览权重
	ArticleLikeHotnessWeight        = 10       // 点赞权重
	ArticleCommentHotnessWeight     = 50       // 评论权重
	ArticleFavoriteHotnessWeight    = 100      // 收藏权重
	ArticleInitialHotness           = 1000000  // 新发布文章初始权重
	ArticleInitialHotnessDecayHours = 180 * 24 // 180天初始权重衰减为1
)

var ArticleHotnessDecayFactor float64 // 时间衰减因子权重

func init() {
	ArticleHotnessDecayFactor = -math.Log(1.0/ArticleInitialHotness) / ArticleInitialHotnessDecayHours // 时间衰减因子
}

type Article struct {
	ID            int64             `gorm:"primaryKey;type:bigint unsigned" json:"id"`
	UserID        int64             `gorm:"type:bigint unsigned;not null" json:"userId,omitempty"`
	Title         string            `gorm:"size:255;not null" json:"title,omitempty"`
	Content       string            `gorm:"not null" json:"content,omitempty"`
	TextContent   string            `gorm:"not null" json:"textContent,omitempty"`
	Preview       string            `gorm:"size255;not null" json:"preview,omitempty"`
	ViewCount     int64             `gorm:"->;type:bigint unsigned;not null;default:0" json:"viewCount,omitempty"`
	LikeCount     int64             `gorm:"type:bigint unsigned;not null;default:0" json:"likeCount,omitempty"`
	CommentCount  int64             `gorm:"type:bigint unsigned;not null;default:0" json:"commentCount,omitempty"`
	FavoriteCount int64             `gorm:"type:bigint unsigned;not null;default:0" json:"favoriteCount,omitempty"`
	Hotness       int64             `gorm:"index:,sort:desc;type:bigint unsigned;not null;default:0" json:"hotness,omitempty"`
	CreatedAt     *time.Time        `gorm:"index;not null" json:"createdAt,omitempty"`
	UpdatedAt     *time.Time        `gorm:"not null" json:"updatedAt,omitempty"`
	DeletedAt     *gorm.DeletedAt   `gorm:"index" json:"deletedAt,omitempty"`
	HasLiked      bool              `gorm:"-" json:"hasLiked,omitempty"`
	HasFavorited  bool              `gorm:"-" json:"hasFavorited,omitempty"`
	Highlights    *map[string][]any `gorm:"-" json:"highlights,omitempty"`
	User          *User             `json:"user,omitempty"`
	Tags          Tags              `gorm:"many2many:article_tags" json:"tags,omitempty"`
	Likes         Likes             `json:"likes,omitempty"`
	Favorites     Favorites         `json:"favorites,omitempty"`
	Comments      Comments          `json:"comments,omitempty"`
}

type Articles []*Article

func (article *Article) BeforeSave(tx *gorm.DB) (err error) {
	var buf bytes.Buffer
	if err = goldmark.Convert([]byte(article.Content), &buf); err != nil {
		return
	}
	article.TextContent = string(bluemonday.StrictPolicy().SanitizeBytes(buf.Bytes()))

	return
}

func (article *Article) Load(preloads map[string][]any, missing ...bool) {
	LoadRelations(app.DB, article, preloads, missing...)
}

func (articles Articles) Load(preloads map[string][]any, missing ...bool) {
	LoadRelations(app.DB, articles, preloads, missing...)
}

func (article *Article) ParseInclude(relations any) {
	availableRelations := map[string][]any{
		"User": {func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Name", "Avatar")
		}},
		"Tags": {func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Name")
		}},
	}

	article.Load(FilterRelations(relations, availableRelations))
}

func (articles Articles) ParseInclude(relations any) {
	availableRelations := map[string][]any{
		"User": {func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Name", "Avatar")
		}},
		"Tags": {func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Name")
		}},
	}

	articles.Load(FilterRelations(relations, availableRelations))
}

func (article *Article) BeforeCreate(_ *gorm.DB) (err error) {
	article.Hotness = article.CalcHotness()
	return
}

func (article *Article) AfterFind(_ *gorm.DB) (err error) {
	field := article.getViewCountIncrField()
	article.ViewCount += field.Current(article.ID)

	return
}

func (article *Article) IncrViewCount(value int64) {
	field := article.getViewCountIncrField()
	field.Incr(article.ID, value)
	article.ViewCount += value
}

func (article *Article) WithUserHasLiked(userId int64) {
	if userId > 0 {
		article.HasLiked = app.DB.Model(&article).Where("user_id = ?", userId).Association("Likes").Count() > 0
	}
}

func (article *Article) WithUserHasFavorited(userId int64) {
	if userId > 0 {
		article.HasFavorited = app.DB.Model(&article).Where("user_id = ?", userId).Association("Favorites").Count() > 0
	}
}

func (article *Article) CalcHotness() int64 {
	date := article.CreatedAt
	if date == nil {
		now := time.Now()
		date = &now
	}

	elapsedHours := time.Since(*date).Hours()
	timeDecay := math.Exp(-ArticleHotnessDecayFactor * elapsedHours)

	return ArticleViewHotnessWeight*article.ViewCount +
		ArticleLikeHotnessWeight*article.LikeCount +
		ArticleCommentHotnessWeight*article.CommentCount +
		ArticleFavoriteHotnessWeight*article.FavoriteCount +
		int64(ArticleInitialHotness*timeDecay)
}

func (article *Article) UpdateHotness() int64 {
	article.Hotness = article.CalcHotness()

	return app.DB.Model(&article).UpdateColumn("hotness", article.Hotness).RowsAffected
}

func (article *Article) Url() *url.URL {
	parsed, _ := url.Parse(fmt.Sprintf("%s/articles/%d", config.App.SiteUrl, article.ID))

	return parsed
}

func (article *Article) getViewCountIncrField() *RedisIncrField {
	return NewRedisIncrField("articles", "view_count", time.Now())
}
