package middleware

import (
	"blog/app"
	h "blog/app/http"
	"blog/app/http/responses"
	"blog/app/validation"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		defer func(start time.Time) {
			if err := recover(); err != nil {
				switch v := err.(type) {
				case *validation.Error:
					responses.Json(c, nil, false, v.Message, http.StatusUnprocessableEntity, responses.CodeParameterError, v.Errors)
					return
				case *h.AuthenticationError:
					responses.Json(c, nil, false, "授权失败，请先登录", http.StatusUnauthorized, responses.CodeUnauthorized)
					return
				case *h.Error:
					responses.Json(c, nil, false, v.Message, v.StatusCode, v.Code)
					if v.Message != "请求参数格式错误" {
						return
					}
				case *mysql.MySQLError:
					responses.Json(c, nil, false, "服务异常", http.StatusInternalServerError, responses.CodeSystemError)
				case string:
					responses.Json(c, nil, false, v, http.StatusBadRequest, responses.CodeBadRequest)
					return
				case error:
					responses.Json(c, nil, false, v.Error(), http.StatusBadRequest, responses.CodeBadRequest)
				default:
					responses.Json(c, nil, false, "未知错误", http.StatusInternalServerError, responses.CodeSystemError)
				}

				app.Logger.With(
					zap.Int64("userId", c.GetInt64("userId")),
					zap.Int("status", c.Writer.Status()),
					zap.Duration("duration", time.Since(start)),
					zap.String("ip", c.ClientIP()),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("query", c.Request.URL.RawQuery),
					zap.String("user-agent", c.Request.UserAgent()),
				).Error(fmt.Sprintf("%T: %v", err, err))
			}
		}(start)

		c.Next()
	}
}
