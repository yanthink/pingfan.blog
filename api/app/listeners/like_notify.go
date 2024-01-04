package listeners

import (
	"blog/app/events"
	"blog/app/notifications"
)

type LikeNotify struct {
	listener
}

func (*LikeNotify) Handle(event any) (err error) {
	e, ok := event.(*events.ArticleLiked)
	if !ok {
		return
	}

	if e.Like.DeletedAt.Valid || e.Like.UpdatedAt != nil && !e.Like.CreatedAt.Equal(*e.Like.UpdatedAt) {
		return
	}

	_ = notifications.Send([]notifications.Notifiable{
		&notifications.ArticleLike{Like: e.Like, Article: e.Article},
	})

	return
}
