package bootstrap

import (
	"blog/app/events"
	"blog/app/listeners"
)

func SetupEvent() {
	events.On((*events.ArticleChanged)(nil), &listeners.ProcessArticleImages{})
	events.On((*events.ArticleLiked)(nil), &listeners.LikeNotify{})
	events.On((*events.Commented)(nil), &listeners.ProcessCommentImages{}, &listeners.CommentNotify{})
}
