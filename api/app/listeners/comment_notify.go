package listeners

import (
	"blog/app/events"
	"blog/app/notifications"
)

type CommentNotify struct {
	listener
}

func (*CommentNotify) Handle(event any) (err error) {
	e, ok := event.(*events.Commented)
	if !ok {
		return nil
	}

	var notifiables []notifications.Notifiable

	if e.Comment.ParentID == 0 {
		notifiables = append(notifiables, &notifications.ArticleComment{Comment: e.Comment})
	} else {
		notifiables = append(
			notifiables,
			&notifications.CommentReply{Comment: e.Comment},
			&notifications.CommentHasNewReply{Comment: e.Comment},
			&notifications.ArticleHasNewReply{Comment: e.Comment})
	}

	_ = notifications.Send(notifiables)

	return
}
