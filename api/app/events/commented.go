package events

import "blog/app/models"

type Commented struct {
	Comment *models.Comment
}
