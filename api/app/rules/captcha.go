package rules

import (
	"blog/app/captcha"
	"github.com/go-playground/validator/v10"
	"reflect"
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

	codeField := fl.Field()
	for codeField.Kind() == reflect.Ptr {
		codeField = codeField.Elem()
	}
	code := fl.Field().String()

	accountField := fl.Top().FieldByName(accountFieldName)
	for accountField.Kind() == reflect.Ptr {
		accountField = accountField.Elem()
	}
	account := accountField.String()

	hashField := fl.Top().FieldByName(hashFieldName)
	for hashField.Kind() == reflect.Ptr {
		hashField = hashField.Elem()
	}
	hash := hashField.String()

	return captcha.New().SetType(CaptchaType).Check(account, code, hash)
}
