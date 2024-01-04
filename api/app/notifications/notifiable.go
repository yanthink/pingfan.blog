package notifications

import "blog/app/models"

type Notifiable interface {
	Setup() error
	Users() models.Users
	FromUser() *models.User
	Subject() string
	Message() string
	ToDatabase() (*map[string]any, bool)
	ToMail() bool
	ToWebsocket() bool
}
