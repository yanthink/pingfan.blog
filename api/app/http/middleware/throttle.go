package middleware

import (
	"blog/app"
	h "blog/app/http"
	"blog/app/http/responses"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"net/http"
	"strconv"
	"time"
)

func Throttle(limit redis_rate.Limit, getKeyFn ...func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()

		if len(getKeyFn) > 0 {
			key = getKeyFn[0](c)
		}

		res, err := app.Limiter.Allow(context.Background(), key, limit)

		if err != nil {
			c.Abort()
			panic(err)
		}

		c.Header("RateLimit-Remaining", strconv.Itoa(res.Remaining))

		if res.Allowed == 0 {
			c.Abort()

			seconds := int(res.RetryAfter / time.Second)
			panic(&h.Error{
				Code:       responses.CodeTooManyAttemptsCode,
				StatusCode: http.StatusTooManyRequests,
				Message:    fmt.Sprintf("请求太频繁，请%d秒后再尝试", seconds),
			})
		}

		c.Next()
	}
}
