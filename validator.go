package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// Result struct definition for a Result object.
type Result struct {
	List any `validate:"dive"`
}

// Validate a new instance of the validator library.
var Validate = validator.New()

// RegExpNumberString a regular expression to match number strings
// with optional thousands separators and decimal portions.
var RegExpNumberString = regexp.MustCompile(`^-?([1-9]{1}\d{0,2}(\,\d{3})*(\.\d+)?|[1-9]{1}\d*(\.\d+)?|0(\.\d+)?|(\.\d+)?)$`)

func init() {
	Validate.RegisterValidation("number_string", ValidateNumberString)
}

// ValidateNumberString validates a given number string by checking whether
// it matches the RegExpNumberString regular expression pattern.
func ValidateNumberString(fl validator.FieldLevel) bool {
	return RegExpNumberString.MatchString(fl.Field().String())
}
