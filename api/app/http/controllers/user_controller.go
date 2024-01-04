package controllers

import (
	"blog/app/helpers"
	r "blog/app/http/requests/user"
	"blog/app/http/responses"
	"blog/app/models"
	"blog/app/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

type userController struct {
}

func (*userController) User(c *gin.Context) {
	responses.Json(c, services.User.GetAuthUser(c))
}

func (*userController) UpdateProfile(c *gin.Context) {
	req := r.UpdateProfileValidate(c)

	var meta *map[string]any

	if req.Meta != nil {
		m := helpers.StructToMap(req.Meta)
		user := services.User.GetAuthUser(c)

		if user.Meta != nil {
			m = helpers.Merge(*user.Meta, m)
		}

		meta = &m
	}

	responses.Json(c, services.User.Update(c.GetInt64("userId"), &models.User{
		Name:   req.Name,
		Email:  req.Email,
		Avatar: req.Avatar,
		Meta:   meta,
	}))
}

func (*userController) UpdateMeta(c *gin.Context) {
	user := services.User.GetAuthUser(c)

	req := r.UpdateMetaValidate(c)
	meta := helpers.StructToMap(req.Meta)
	if user.Meta != nil {
		meta = helpers.Merge(*user.Meta, meta)
	}

	responses.Json(c, services.User.Update(c.GetInt64("userId"), &models.User{
		Meta: &meta,
	}))
}

func (*userController) UpdatePassword(c *gin.Context) {
	req := r.UpdatePasswordValidate(c)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	responses.Json(c, services.User.Update(c.GetInt64("userId"), &models.User{
		Password: string(hashedPassword),
	}))
}

func (*userController) UnreadNotificationCount(c *gin.Context) {
	responses.Json(c, services.Notification.UnreadCount(c.GetInt64("userId")))
}

func (*userController) Notifications(c *gin.Context) {
	notifications, count := services.Notification.Paginate(r.NotificationsValidate(c))

	responses.Json(c, notifications, count)
}

func (*userController) NotificationsMarkAsRead(c *gin.Context) {
	responses.Json(c, gin.H{"rows": services.Notification.MarkAsRead(c.GetInt64("userId"))})
}

func (*userController) Favorites(c *gin.Context) {
	favorites, count := services.Favorite.Paginate(r.FavoritesValidate(c))

	responses.Json(c, favorites, count)
}

func (*userController) Comments(c *gin.Context) {
	comments, count := services.Comment.Paginate(r.CommentsValidate(c))

	responses.Json(c, comments, count)
}

func (*userController) Likes(c *gin.Context) {
	likes, count := services.Like.Paginate(r.LikesValidate(c))

	responses.Json(c, likes, count)
}

func (*userController) Upvotes(c *gin.Context) {
	upvotes, count := services.Upvote.Paginate(r.UpvotesValidate(c))

	responses.Json(c, upvotes, count)
}

func (*userController) Index(c *gin.Context) {
	users, count := services.User.Paginate(r.IndexValidate(c))

	responses.Json(c, users, count)
}

func (*userController) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id < 1 {
		responses.Json(c, nil, false, http.StatusNotFound, responses.CodeNotFound)
		return
	}

	req := r.UpdateValidate(c)

	services.User.Update(id, &models.User{
		Status: req.Status,
	})
}
