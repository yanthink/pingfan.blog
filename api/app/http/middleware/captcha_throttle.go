package middleware

import (
	"blog/app"
	h "blog/app/http"
	"blog/app/http/responses"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
)

func CaptchaThrottle(t string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("rate_%s:%s", t, c.ClientIP())

		defer func() {
			if r := recover(); r != nil {
				if e, ok := r.(*h.Error); !ok || e.Code != responses.CodeTooManyAttemptsCode {
					_ = app.Limiter.Reset(context.Background(), key)
				}

				panic(r)
			}
		}()

		Throttle(redis_rate.PerMinute(1), func(*gin.Context) string {
			return key
		})(c)
	}
}
