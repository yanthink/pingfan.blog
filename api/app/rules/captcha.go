package rules

import (
	"blog/app/captcha"
	"github.com/go-playground/validator/v10"
	"strings"
)

// ValidateCaptcha 验证码 eg: binding:"required,captcha=Email CodeKey type:email"
func ValidateCaptcha(fl validator.FieldLevel) bool {
	accountFieldName := "Email"
	hashFieldName := "CodeKey"
	CaptchaType := captcha.Email

	param := fl.Param()
	if param != "" {
		for i, s := range strings.Fields(fl.Param()) {
			switch i {
			case 0:
				accountFieldName = s
			case 1:
				hashFieldName = s
			case 2:
				if strings.HasPrefix(s, "type:") {
					CaptchaType = s[5:]
				}
			}
		}
	}

	code := fl.Field().String()
	account := fl.Top().FieldByName(accountFieldName).String()
	hash := fl.Top().FieldByName(hashFieldName).String()

	return captcha.New().SetType(CaptchaType).Check(account, code, hash)
}
