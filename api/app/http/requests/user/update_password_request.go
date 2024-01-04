package user

import (
	h "blog/app/http"
	"blog/app/http/responses"
	"blog/app/services"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UpdatePasswordRequest struct {
	OldPassword          string `form:"oldPassword" json:"oldPassword" label:"旧密码"`
	Password             string `form:"password" json:"password" label:"新密码" binding:"required,nefield=OldPassword,min=6"`
	PasswordConfirmation string `form:"passwordConfirmation" json:"passwordConfirmation" label:"确认密码" binding:"required,eqfield=Password"`
}

func UpdatePasswordValidate(c *gin.Context) (r *UpdatePasswordRequest) {
	r = &UpdatePasswordRequest{}
	validation.Validate(c, r)

	user := services.User.GetAuthUser(c)

	if user.HasPassword {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.OldPassword)); err != nil {
			panic(&h.Error{Code: responses.CodeParameterError, StatusCode: http.StatusUnprocessableEntity, Message: "旧密码不正确"})
		}
	}

	return
}
