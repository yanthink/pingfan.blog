package middleware

import (
	h "blog/app/http"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"strings"
)

func Auth(jwtKey []byte, abort ...bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if len(abort) == 0 || abort[0] {
					c.Abort()
					panic(err)
				}
			}
		}()

		token := strings.TrimPrefix(c.DefaultQuery("token", c.GetHeader("Authorization")), "Bearer ")

		if token == "" {
			panic(&h.AuthenticationError{})
		}

		t, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
			return jwtKey, nil
		})

		if err != nil {
			panic(&h.AuthenticationError{})
		}

		var userId int64

		if claims, ok := t.Claims.(*jwt.RegisteredClaims); ok && t.Valid {
			userId, _ = strconv.ParseInt(claims.Subject, 10, 64)
		}

		if userId == 0 {
			panic(&h.AuthenticationError{})
		}

		c.Set("userId", userId)
	}
}
