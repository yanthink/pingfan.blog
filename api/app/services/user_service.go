package services

import (
	"blog/app"
	"blog/app/cache"
	"blog/app/filters"
	h "blog/app/http"
	"blog/app/http/responses"
	"blog/app/mini_program"
	"blog/app/models"
	"blog/app/pagination"
	"blog/app/resource"
	"blog/app/storage"
	"blog/app/websocket"
	"blog/config"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"strconv"
	"time"
)

type userService struct {
}

func (*userService) Login(name, password string) (signedToken string, user *models.User) {
	app.DB.Where("name = ?", name).First(&user)
	abortIf(user.ID == 0 || user.Password == "", "用户名或密码不正确", http.StatusUnprocessableEntity, responses.CodeLoginError)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	abortIf(err != nil, "用户名或密码不正确", http.StatusUnprocessableEntity, responses.CodeLoginError)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(user.ID, 10),
		IssuedAt:  jwt.NewNumericDate(time.Now()), // 签发时间
		NotBefore: jwt.NewNumericDate(time.Now()), // 生效时间
	})

	if signedToken, err = token.SignedString(config.Jwt.Key); err != nil {
		panic(err)
	}

	return
}

func (*userService) WxLogin(code string) (signedToken string, user *models.User) {
	session := mini_program.Wx.Code2Session(code, "")

	err := app.DB.Transaction(func(tx *gorm.DB) (err error) {
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where(models.User{Openid: &session.Openid}).
			FirstOrCreate(&user).Error
		return
	})

	if err != nil {
		panic(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(user.ID, 10),
		IssuedAt:  jwt.NewNumericDate(time.Now()), // 签发时间
		NotBefore: jwt.NewNumericDate(time.Now()), // 生效时间
	})

	if signedToken, err = token.SignedString(config.Jwt.Key); err != nil {
		panic(err)
	}

	return
}

var wxLoginQRCodeSnowflakeNode, _ = snowflake.NewNode(config.Snowflake.Node)

func (*userService) GetWxLoginQRCode() (signedToken, img string, expiration time.Duration) {
	uuid := fmt.Sprintf(fmt.Sprintf("uid_%s", wxLoginQRCodeSnowflakeNode.Generate().String()))

	envVersion := "release"
	if config.App.Env == "local" {
		envVersion = "trial"
	}

	buffer := mini_program.Wx.GetUnlimitedQRCode(&mini_program.GetUnlimitedQRCodeParams{
		Scene:      uuid,
		Page:       "pages/scan/login",
		CheckPath:  false,
		EnvVersion: envVersion,
		Width:      328,
		IsHyaline:  true,
	})

	expiration = 5 * time.Minute
	if config.App.Env == "local" {
		expiration = 30 * time.Second
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   uuid,
		IssuedAt:  jwt.NewNumericDate(time.Now()), // 签发时间
		NotBefore: jwt.NewNumericDate(time.Now()), // 生效时间
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
	})

	signedToken, _ = token.SignedString(config.Jwt.Key)
	img = base64.StdEncoding.EncodeToString(buffer.Bytes())

	strCache := cache.New[string]()
	cacheKey := fmt.Sprintf("wx_scan_login:%s", uuid)
	_ = strCache.Put(context.Background(), cacheKey, "1", expiration)

	return
}

func (*userService) WxScanLogin(user *models.User, uuid string) {
	strCache := cache.New[string]()
	cacheKey := fmt.Sprintf("wx_scan_login:%s", uuid)

	if v, _ := strCache.Pull(context.Background(), cacheKey); v != "1" {
		abort("小程序码无效或已过期！")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(user.ID, 10),
		IssuedAt:  jwt.NewNumericDate(time.Now()), // 签发时间
		NotBefore: jwt.NewNumericDate(time.Now()), // 生效时间
	})

	signedToken, err := token.SignedString(config.Jwt.Key)
	if err != nil {
		panic(err)
	}

	websocket.SendToTempUser(&websocket.TempUserMessage{
		TempUserID: uuid,
		Response: &websocket.Response{
			Event: "WxScanLoginSuccess",
			Data: map[string]any{
				"token": signedToken,
				"user":  user,
			},
		},
	})
}

func (*userService) GetById(id int64) (user *models.User, err error) {
	err = app.DB.First(&user, id).Error

	return
}

func (*userService) GetByEmail(email string) (user *models.User, err error) {
	err = app.DB.Where("email = ?", email).First(&user).Error

	return
}

func (*userService) GetByName(name string) (user *models.User, err error) {
	err = app.DB.Where("name = ?", name).First(&user).Error

	return
}

func (s *userService) GetAuthUser(c *gin.Context) (user *models.User) {
	if v, ok := c.Get("user"); ok {
		if user, ok = v.(*models.User); ok && user != nil {
			return
		}
	}

	if user, _ = s.GetById(c.GetInt64("userId")); user.ID == 0 || *user.Status == 1 {
		panic(&h.AuthenticationError{})
	}

	c.Set("user", user)

	return
}

func (s *userService) CheckAuthIsAdmin(c *gin.Context) (user *models.User) {
	user = s.GetAuthUser(c)
	abortIf(user == nil || user.Role != models.UserRoleManage, "无权限操作", responses.CodeAccessDenied, http.StatusForbidden)

	return
}

func (s *userService) Update(id int64, user *models.User) *models.User {
	user.ID = id

	rType := resource.Avatar
	if user.Avatar != "" && rType.IsUploadPath(storage.Disk().ParsePath(user.Avatar)) {
		user.Avatar, _ = Resource.CopyToStorePath(user.Avatar, rType)
	}

	app.DB.Model(user).Omit(clause.Associations, "Openid").Updates(user)

	return user
}

func (*userService) Paginate(paginator pagination.Pager) (users models.Users, count int64) {
	tx := app.DB.
		Model(users).
		Order("id DESC").
		Scopes(filters.New(&filters.UserFilter{}, paginator))

	var err error

	if _, count, err = paginator.Paginate(tx, &users); err != nil {
		panic(err)
	}

	return
}
