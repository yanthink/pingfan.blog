package services

import (
	h "blog/app/http"
	"blog/app/http/responses"
	"net/http"
)

var (
	Article      *articleService
	User         *userService
	Resource     *resourceService
	Tag          *tagService
	Comment      *commentService
	Notification *notificationService
	Favorite     *favoriteService
	Like         *likeService
	Upvote       *upvoteService
	Captcha      *captchaService
	Zincsearch   *zincsearchService
)

func abort(message string, params ...any) {
	code := responses.CodeBadRequest
	statusCode := http.StatusBadRequest

	for _, param := range params {
		switch v := param.(type) {
		case responses.Code:
			code = v
		case int:
			statusCode = v
		}
	}

	panic(&h.Error{
		Code:       code,
		StatusCode: statusCode,
		Message:    message,
	})
}

func abortIf(b bool, message string, params ...any) {
	if b {
		abort(message, params...)
	}
}
