package login

import (
	"blog/app/models"
	"blog/app/services"
	"blog/app/validation"
	"github.com/gin-gonic/gin"
)

type AccountRequest struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (r *AccountRequest) Login() (string, *models.User) {
	return services.User.Login(r.Name, r.Password)
}

func AccountValidate(c *gin.Context) (r *AccountRequest) {
	r = &AccountRequest{}
	validation.Validate(c, r)

	return
}
