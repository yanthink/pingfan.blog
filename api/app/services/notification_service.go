package services

import (
	"blog/app"
	"blog/app/filters"
	"blog/app/models"
	"blog/app/pagination"
	"time"
)

type notificationService struct {
}

func (*notificationService) Paginate(paginator pagination.Pager) (notifications models.Notifications, count int64) {
	tx := app.DB.
		Model(notifications).
		Order("id DESC").
		Scopes(filters.New(&filters.NotificationFilter{}, paginator))

	var err error

	if _, count, err = paginator.Paginate(tx, &notifications); err != nil {
		panic(err)
	}

	return
}

func (*notificationService) UnreadCount(userId int64) (count int64) {
	app.DB.Model(&models.Notification{}).Where("user_id = ?", userId).Where("read_at IS NULL").Count(&count)

	return
}

func (*notificationService) MarkAsRead(userId int64) int64 {
	return app.DB.
		Model(&models.Notification{}).
		Where("user_id = ?", userId).
		Where("read_at IS NULL").
		Update("read_at", time.Now()).
		RowsAffected
}
