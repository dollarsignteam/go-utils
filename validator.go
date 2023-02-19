package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

var RegExpNumberString = regexp.MustCompile(`^-?([1-9]{1}\d{0,2}(\,\d{3})*(\.\d+)?|[1-9]{1}\d*(\.\d+)?|0(\.\d+)?|(\.\d+)?)$`)

func init() {
	Validate.RegisterValidation("number_string", ValidateNumberString)
}

func ValidateNumberString(fl validator.FieldLevel) bool {
	return RegExpNumberString.MatchString(fl.Field().String())
}
