package utils

import (
	"reflect"
	"regexp"
	"strings"

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
	Validate.RegisterTagNameFunc(GetJSONTagName)
	_ = Validate.RegisterValidation("number_string", ValidateNumberString)
}

// ValidateNumberString validates a given number string by checking whether
// it matches the RegExpNumberString regular expression pattern.
func ValidateNumberString(fl validator.FieldLevel) bool {
	return RegExpNumberString.MatchString(fl.Field().String())
}

// GetJSONTagName returns the name of the JSON tag associated with a given reflect.StructField.
func GetJSONTagName(field reflect.StructField) string {
	tagName := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
	if tagName == "-" {
		return ""
	}
	return tagName
}

// ValidateStruct validates the fields of a struct using the validator library.
// if validation fails, it calls the ParseValidationError() function
// to convert ValidationErrors into a ValidationError error type.
func ValidateStruct(s any) error {
	if err := Validate.Struct(s); err != nil {
		return ParseValidationError(err)
	}
	return nil
}
