package middleware

import (
	h "blog/app/http"
	"blog/app/services"
	"github.com/gin-gonic/gin"
)

func CheckUserStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := services.User.GetAuthUser(c)
		if *user.Status == 1 {
			panic(&h.AuthenticationError{})
		}
	}
}
