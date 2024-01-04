package listeners

import (
	"blog/app/events"
	"blog/app/notifications"
)

type UpvoteNotify struct {
	listener
}

func (*UpvoteNotify) Handle(event any) (err error) {
	e, ok := event.(*events.CommentUpvoted)
	if !ok {
		return
	}

	if e.Upvote.ID == 0 || e.Upvote.UpdatedAt != nil && !e.Upvote.CreatedAt.Equal(*e.Upvote.UpdatedAt) {
		return
	}

	_ = notifications.Send([]notifications.Notifiable{
		&notifications.CommentUpvote{Upvote: e.Upvote, Comment: e.Comment},
	})

	return
}
