package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
)

func ActionThrottle(action string) gin.HandlerFunc {
	return Throttle(redis_rate.PerMinute(4), func(c *gin.Context) string {
		return fmt.Sprintf("rate_%s:%s:%d", action, c.Param("id"), c.GetInt64("userId"))
	})
}
