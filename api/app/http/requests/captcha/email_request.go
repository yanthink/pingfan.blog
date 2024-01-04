package captcha

import (
	h "blog/app/http"
	"blog/app/http/responses"
	"blog/app/services"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

type EmailRequest struct {
	Email       string `form:"email" json:"email" binding:"required,email"`
	CheckUnique bool   `form:"checkUnique" json:"checkUnique"`
}

func EmailValidate(c *gin.Context) (r *EmailRequest) {
	r = &EmailRequest{}
	validation.Validate(c, r)

	if r.CheckUnique {
		if _, err := services.User.GetByEmail(r.Email); err == nil {
			panic(&h.Error{Code: responses.CodeParameterError, StatusCode: http.StatusUnprocessableEntity, Message: "邮箱已经存在。"})
		}
	}

	return
}
