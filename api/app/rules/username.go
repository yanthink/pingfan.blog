package rules

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidateUsername(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	// go 正则不支持负前瞻
	if value[0] == '_' || value[len(value)-1] == '_' {
		return false
	}

	match, _ := regexp.MatchString(`^[a-zA-Z0-9_\p{Han}]{2,10}$`, value)

	return match
}
