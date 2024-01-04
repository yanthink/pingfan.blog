package rules

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidatePhone(fl validator.FieldLevel) bool {
	match, _ := regexp.MatchString(`^1[3456789]\d{9}$`, fl.Field().String())

	return match
}
