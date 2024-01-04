package routes

import (
	"blog/app/http/controllers"
	"blog/app/http/middleware"
	"blog/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
)

func RegisterApiRouter(r *gin.Engine) {
	r.Static("storage", "storage/app/public")

	api := r.Group("api").
		Use(middleware.Throttle(redis_rate.PerMinute(60))).
		Use(middleware.Auth(config.Jwt.Key, false))
	{
		api.POST("login", controllers.Login.Account)
		api.POST("login/wx", controllers.Login.Wx)
		api.POST("login/wx_qrcode", controllers.Login.WxQRCode)
		api.POST("login/wx_scan", controllers.Login.WxScan)

		api.GET("articles", controllers.Article.Index)
		api.GET("articles/search", controllers.Article.Search)
		api.GET("articles/cursor_paginate", controllers.Article.CursorPaginate)
		api.GET("articles/:id", controllers.Article.Show)

		api.GET("comments", controllers.Comment.Index)
		api.GET("comments/cursor_paginate", controllers.Comment.CursorPaginate)
		api.GET("comments/:id", controllers.Comment.Show)

		api.GET("tags", controllers.Tag.Index)
	}

	authorized := r.Group("api").
		Use(middleware.Throttle(redis_rate.PerMinute(60))).
		Use(middleware.Auth(config.Jwt.Key)).
		Use(middleware.CheckUserStatus())
	{
		authorized.GET("user", controllers.User.User)
		authorized.PUT("user", controllers.User.UpdateProfile)
		authorized.PUT("user/meta", controllers.User.UpdateMeta)
		authorized.PUT("user/password", controllers.User.UpdatePassword)
		authorized.GET("user/unread_notification_count", controllers.User.UnreadNotificationCount)
		authorized.GET("user/notifications", controllers.User.Notifications)
		authorized.POST("user/notifications_mark_as_read", controllers.User.NotificationsMarkAsRead)
		authorized.GET("user/favorites", controllers.User.Favorites)
		authorized.GET("user/comments", controllers.User.Comments)
		authorized.GET("user/likes", controllers.User.Likes)
		authorized.GET("user/upvotes", controllers.User.Upvotes)

		authorized.GET("users", controllers.User.Index)
		authorized.PUT("users/:id", controllers.User.Update)

		authorized.POST("articles", controllers.Article.Store)
		authorized.PUT("articles/:id", controllers.Article.Update)
		authorized.POST("articles/:id/like", middleware.ActionThrottle("like"), controllers.Article.Like)
		authorized.POST("articles/:id/favorite", middleware.ActionThrottle("favorite"), controllers.Article.Favorite)

		authorized.POST("comments", controllers.Comment.Store)
		authorized.PUT("comments/:id", controllers.Comment.Update)
		authorized.POST("comments/:id/upvote", middleware.ActionThrottle("upvote"), controllers.Comment.Upvote)

		authorized.POST("tags", controllers.Tag.Index)
		authorized.PUT("tags/:id", controllers.Tag.Update)

		authorized.POST("resources/upload", controllers.Resource.Upload)
	}

	captcha := r.Group("api/captcha").Use(middleware.CaptchaThrottle("email"))
	{
		captcha.POST("email", controllers.Captcha.Email)
	}
}
