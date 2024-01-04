package models

import (
	"blog/app"
	"blog/app/helpers"
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	ID          int64           `gorm:"primaryKey;type:bigint unsigned" json:"id"`
	UserID      int64           `gorm:"index;type:bigint unsigned;not null" json:"userId,omitempty"`
	ArticleID   int64           `gorm:"index;type:bigint unsigned;not null" json:"articleId,omitempty"`
	ParentID    int64           `gorm:"type:bigint unsigned;not null" json:"parentId,omitempty"`
	CommentID   int64           `gorm:"index;type:bigint unsigned;not null;comment:>=2级的评论设置成1级评论ID" json:"commentId,omitempty"`
	Content     string          `gorm:"size:2048;not null" json:"content,omitempty"`
	UpvoteCount int64           `gorm:"index:,sort:desc;type:bigint unsigned;not null;default:0" json:"upvoteCount,omitempty"`
	ReplyCount  int64           `gorm:"index:,sort:desc;type:bigint unsigned;not null;default:0" json:"replyCount,omitempty"`
	CreatedAt   *time.Time      `gorm:"not null" json:"createdAt,omitempty"`
	UpdatedAt   *time.Time      `gorm:"not null" json:"updatedAt,omitempty"`
	DeletedAt   *gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	HasUpvoted  bool            `gorm:"-" json:"hasUpvoted,omitempty"`
	User        *User           `json:"user,omitempty"`
	Article     *Article        `json:"article,omitempty"`
	Replies     Comments        `json:"replies,omitempty"`
	Parent      *Comment        `json:"parent,omitempty"`
	Comment     *Comment        `json:"comment,omitempty"`
	Upvotes     Upvotes         `json:"upvotes,omitempty"`
}

type Comments []*Comment

func (comment *Comment) Load(preloads map[string][]any, missing ...bool) {
	LoadRelations(app.DB, comment, preloads, missing...)
}

func (comments Comments) Load(preloads map[string][]any, missing ...bool) {
	LoadRelations(app.DB, comments, preloads, missing...)
}

func (comment *Comment) ParseInclude(relations any) {
	availableRelations := map[string][]any{
		"User": {func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Name", "Avatar")
		}},
		"Article": {func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Title")
		}},
	}

	comment.Load(FilterRelations(relations, availableRelations))
}

func (comments Comments) ParseInclude(relations any) {
	availableRelations := map[string][]any{
		"User": {func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Name", "Avatar")
		}},
		"Article": {func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Title")
		}},
		"Replies": {func(db *gorm.DB) *gorm.DB {
			return db.
				// Session(&gorm.Session{NewDB: true}). // 这里不能使用NewDB，否则后续的关系无法预/延时加载
				Table("(?) as c", db.Session(&gorm.Session{NewDB: true}).Model(&Comment{}).Select("*, ROW_NUMBER() OVER(PARTITION BY comment_id) as `rank`")).
				Where("c.rank <= 3")
		}},
		"Replies.User": {func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Name", "Avatar")
		}},
		"Parent.User": {func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Name", "Avatar")
		}},
	}

	if len(comments) > 0 && comments[0].CommentID > 0 {
		delete(availableRelations, "Replies")
		delete(availableRelations, "Replies.User")
	}

	comments.Load(FilterRelations(relations, availableRelations))
}

func (comment *Comment) WithUserHasUpvoted(userId int64) {
	if userId > 0 {
		comment.HasUpvoted = app.DB.Model(&comment).Where("user_id = ?", userId).Association("Upvotes").Count() > 0
	}
}

func (comments Comments) WithUserHasUpvoted(userId int64) {
	if userId > 0 && len(comments) > 0 {
		var upvotes Upvotes
		_ = app.DB.Model(comments).Where("user_id = ?", userId).Select("CommentID").Association("Upvotes").Find(&upvotes)
		upvoteMap := helpers.KeyBy(upvotes, func(upvote *Upvote) int64 {
			return upvote.CommentID
		})

		for _, comment := range comments {
			if _, ok := upvoteMap[comment.ID]; ok {
				comment.HasUpvoted = true
			}
		}
	}
}

func (comments Comments) Flat() Comments {
	var flatComments Comments

	for _, comment := range comments {
		flatComments = append(flatComments, comment)

		if comment.Replies != nil {
			for _, reply := range comment.Replies {
				flatComments = append(flatComments, reply)
			}
		}
	}

	return flatComments
}
