package events

import "blog/app/models"

type CommentUpvoted struct {
	Upvote  *models.Upvote
	Comment *models.Comment
}
