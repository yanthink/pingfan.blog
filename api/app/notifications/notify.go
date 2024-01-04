package notifications

import (
	"blog/app"
	"blog/app/helpers"
	"blog/app/mail"
	"blog/app/markdown"
	"blog/app/models"
	"blog/app/websocket"
	"blog/config"
	"bytes"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"go.uber.org/zap"
	"reflect"
)

func Send(notifiables []Notifiable) (err error) {
	defer func() {
		if r := recover(); r != nil {
			app.Logger.Sugar().Error("发送通知失败", r)
		}
	}()

	for _, notifiable := range notifiables {
		if notifiable == nil {
			continue
		}

		// todo 多个 goroutine 同时运行 Setup 加载 model 关系 可能会导致 不同 goroutine 的关系被覆盖，所以这里改成同步运行
		if err = notifiable.Setup(); err != nil {
			return
		}

		users := notifiable.Users()
		if len(users) < 1 {
			return
		}

		subject := notifiable.Subject()
		if subject == "" {
			return
		}

		var buf bytes.Buffer
		_ = goldmark.Convert([]byte(subject), &buf)
		stripTagsSubject := string(bluemonday.StrictPolicy().SanitizeBytes(buf.Bytes()))

		message := notifiable.Message()

		for _, user := range users {
			t := reflect.TypeOf(notifiable)
			for t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			event := t.Name()

			// 新建数据库通知
			if data, ok := notifiable.ToDatabase(); ok {
				var formUserId int64
				if fromUser := notifiable.FromUser(); fromUser != nil {
					formUserId = fromUser.ID
				}

				err = app.DB.Create(&models.Notification{
					UserID:     user.ID,
					FromUserID: formUserId,
					Type:       event,
					Subject:    subject,
					Message:    message,
					Data:       data,
				}).Error

				if err != nil {
					return err
				}
			}

			// websocket 通知
			if notifiable.ToWebsocket() {
				websocket.SendToUser(&websocket.UserMessage{
					UserID: uint64(user.ID),
					Response: &websocket.Response{
						Event: fmt.Sprintf("notifications.%s", event),
						Data: map[string]any{
							"message": stripTagsSubject,
						},
					},
				})
			}
		}

		// 邮件通知
		if notifiable.ToMail() {
			users = helpers.Filter(users, func(_ int, user *models.User) bool {
				if user.Email == "" || user.Meta == nil {
					return false
				}

				emailNotify := (*user.Meta)["emailNotify"].(int64)
				if emailNotify == 0 || emailNotify == 2 && websocket.IsOnline(uint64(user.ID)) {
					return false
				}

				return true
			})

			if len(users) > 0 {
				body, _ := markdown.Parse([]byte(fmt.Sprintf("%s\n\n%s\n\nThanks.\n\n<h3>%s</h3>", subject, message, config.App.Name)))

				err = mail.New().
					AppendTo(users[0].Email).
					AppendCc(helpers.Map(users[1:], func(_ int, user *models.User) string {
						return user.Email
					})...).
					Send(stripTagsSubject, body)

				if err != nil {
					app.Logger.Debug("邮件通知发送失败", zap.Error(err))
				}
			}
		}
	}

	return
}
