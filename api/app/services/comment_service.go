package services

import (
	"blog/app"
	"blog/app/events"
	"blog/app/filters"
	"blog/app/http/responses"
	"blog/app/models"
	"blog/app/pagination"
	"blog/config"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"time"
)

type commentService struct {
}

func (*commentService) Paginate(paginator pagination.Pager) (comments models.Comments, count int64) {
	tx := app.DB.
		Model(comments).
		Omit("UpdatedAt").
		Scopes(filters.New(&filters.CommentFilter{}, paginator))

	var err error

	if _, count, err = paginator.Paginate(tx, &comments); err != nil {
		panic(err)
	}

	return
}

func (*commentService) CursorPaginate(paginator pagination.CursorPager) (comments models.Comments, cursor *pagination.Cursor) {
	tx := app.DB.
		Model(comments).
		Omit("UpdatedAt").
		Scopes(filters.New(&filters.CommentFilter{}, paginator))

	var err error

	if _, cursor, err = paginator.Paginate(tx, &comments); err != nil {
		panic(err)
	}

	return
}

func (*commentService) GetByID(id int64) (comment *models.Comment) {
	abortIf(app.DB.First(&comment, id).RowsAffected == 0, "数据不存在", responses.CodeModelNotFound, http.StatusNotFound)

	return
}

func (s *commentService) Add(comment *models.Comment) *models.Comment {
	if err := app.DB.Create(&comment).Error; err != nil {
		panic(err)
	}

	app.DB.Model(&models.Article{ID: comment.ArticleID}).UpdateColumn("comment_count", gorm.Expr("comment_count + 1"))

	if comment.CommentID > 0 {
		app.DB.Model(&models.Comment{ID: comment.CommentID}).Updates(map[string]any{
			"reply_count": gorm.Expr("reply_count + 1"),
		})
	}

	_ = events.Fire(&events.Commented{Comment: comment})

	return comment
}

func (s *commentService) Update(id int64, comment *models.Comment) *models.Comment {
	comment.ID = id

	if err := app.DB.Model(comment).Omit(clause.Associations).Updates(comment).Error; err != nil {
		panic(err)
	}

	return comment
}

func (s *commentService) Upvote(id, userId int64) *models.Comment {
	lockKey := fmt.Sprintf("%s_upvote_lock:%d", config.Redis.Prefix, userId)
	mutex := app.Redsync.NewMutex(lockKey, redsync.WithTries(1))
	abortIf(mutex.Lock() != nil, "请勿重复操作")
	defer mutex.Unlock()

	comment := s.GetByID(id)
	upvote := models.Upvote{CommentID: comment.ID, UserID: userId}

	if app.DB.Unscoped().Where(&upvote).FirstOrInit(&upvote).RowsAffected == 0 {
		app.DB.Create(&upvote)
	} else {
		if upvote.DeletedAt.Valid { // 恢复
			upvote.DeletedAt = &gorm.DeletedAt{}
		} else { // 删除
			upvote.DeletedAt = &gorm.DeletedAt{
				Time:  time.Now(),
				Valid: true,
			}
		}

		app.DB.Unscoped().Select("DeletedAt").Updates(&upvote)
	}

	comment.HasUpvoted = !upvote.DeletedAt.Valid

	if comment.HasUpvoted {
		comment.UpvoteCount += 1
		app.DB.Model(&comment).UpdateColumn("upvote_count", gorm.Expr("upvote_count + 1"))
	} else {
		comment.UpvoteCount -= 1
		app.DB.Model(&comment).UpdateColumn("upvote_count", gorm.Expr("upvote_count - 1"))
	}

	return comment
}
