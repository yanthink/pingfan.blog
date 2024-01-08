package user

import (
	"blog/app/services"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
)

type UpdateProfileRequest struct {
	Name         *string `form:"name" json:"name" label:"用户名" binding:"required,username" message:"用户名格式不正确"`
	Email        *string `form:"email" json:"email" binding:"omitempty,email"`
	EmailCode    string  `form:"emailCode" json:"emailCode" binding:"omitempty,captcha=Email EmailCodeKey" message:"验证码不正确"`
	EmailCodeKey string  `form:"emailCodeKey" json:"emailCodeKey"`
	Avatar       string  `form:"avatar" json:"avatar" label:"头像" binding:"omitempty,url"`
	Meta         *Meta   `form:"meta" json:"meta"`
}

func UpdateProfileValidate(c *gin.Context) (r *UpdateProfileRequest) {
	r = &UpdateProfileRequest{}
	validation.Validate(c, r)

	user := services.User.GetAuthUser(c)

	if r.Email != nil && (user.Email == nil || *r.Email != *user.Email) {
		if r.EmailCode == "" {
			panic(&validation.Error{
				Message: "邮箱验证码不能为空",
				Errors: map[string][]string{
					"emailCode": {"邮箱验证码不能为空"},
				},
			})
		}
	}

	if user.Name == nil || *r.Name != *user.Name {
		if _, err := services.User.GetByName(*r.Name); err == nil {
			panic(&validation.Error{
				Message: "用户名已经存在",
				Errors: map[string][]string{
					"name": {"用户名已经存在"},
				},
			})
		}
	}

	return
}
