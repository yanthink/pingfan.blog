package http

import "blog/app/http/responses"

type AuthenticationError struct {
}

func (e *AuthenticationError) Error() string {
	return "Unauthorized"
}

type Error struct {
	Code       responses.Code
	StatusCode int
	Message    string
}

func (e *Error) Error() string {
	return e.Message
}
