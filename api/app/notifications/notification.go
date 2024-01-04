package notifications

import "blog/app/models"

type notification struct {
}

func (*notification) Setup() error {
	return nil
}

func (c *notification) Users() models.Users {
	return nil
}

func (c *notification) FromUser() (user *models.User) {
	return
}

func (c *notification) Subject() string {
	return ""
}

func (c *notification) Message() string {
	return ""
}

func (c *notification) ToDatabase() (*map[string]any, bool) {
	return nil, true
}

func (*notification) ToMail() bool {
	return true
}

func (*notification) ToWebsocket() bool {
	return true
}
