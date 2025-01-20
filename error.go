package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Error code constants
const (
	ErrCodeSomethingWentWrong = "SOMETHING_WENT_WRONG"
	ErrCodeBadRequest         = "BAD_REQUEST"
	ErrCodeUnauthorized       = "UNAUTHORIZED"
	ErrCodeNotFound           = "NOT_FOUND"
)

// Error message constants
const (
	ErrMessageSomethingWentWrong = "Something went wrong"
	ErrMessageBadRequest         = "Bad request"
	ErrMessageUnauthorized       = "Unauthorized"
	ErrMessageNotFound           = "Not Found"
)

const errMessageValidationFailed = "Key: '%s', Error: Validation for '%s' failed on the '%s' tag"

// CommonError type for generic errors with status codes and error codes
type CommonError struct {
	StatusCode    int    `json:"statusCode"`   // HTTP status code
	ErrorCode     string `json:"errorCode"`    // Specific error code
	ErrorInstance error  `json:"errorMessage"` // The actual error instance
}

// ValidationError represents an error related to validation, with error details.
type ValidationError struct {
	ErrorMessage string                     `json:"errorMessage"`      // Overall error message
	Details      []ValidationErrorDetail    `json:"details,omitempty"` // Optional list of error details
	Errors       validator.ValidationErrors `json:"-"`                 // The actual validation errors
}

// ValidationErrorDetail represents an individual error detail, with a field, tag, and message.
type ValidationErrorDetail struct {
	Field   string `json:"field" example:"ID"`                                                                              // Field that caused the validation error
	Tag     string `json:"tag" example:"required"`                                                                          // Validation tag that caused the error
	Message string `json:"message" example:"Key: 'Member.ID' Error:Field validation for 'ID' failed on the 'required' tag"` // Full error message
}

// ErrorResponse represents an error response with error code, message, description, and validation errors
type ErrorResponse struct {
	StatusCode       int                     `json:"statusCode" example:"500"`                                     // HTTP status code
	ErrorCode        string                  `json:"errorCode,omitempty" example:"SOMETHING_WENT_WRONG"`           // Specific error code
	ErrorMessage     string                  `json:"errorMessage,omitempty" example:"Oops, something went wrong!"` // Custom error message
	ErrorDescription string                  `json:"errorDescription,omitempty" example:"Something went wrong"`    // The actual error message
	ErrorValidation  []ValidationErrorDetail `json:"errorValidation,omitempty"`                                    // List of validation errors
}

// Error function for CommonError to return the error message
func (e CommonError) Error() string {
	if e.ErrorInstance != nil {
		return e.ErrorInstance.Error()
	}
	return ErrMessageSomethingWentWrong
}

// MarshalJSON function for CommonError to marshal the error as JSON
func (e CommonError) MarshalJSON() ([]byte, error) {
	type Alias CommonError
	return json.Marshal(&struct {
		Alias
		ErrorMessage string `json:"errorMessage"`
	}{
		Alias:        (Alias)(e),
		ErrorMessage: e.Error(),
	})
}

// UnmarshalJSON function for CommonError to unmarshal the error from JSON
func (e *CommonError) UnmarshalJSON(data []byte) error {
	type Alias CommonError
	aux := &struct {
		*Alias
		ErrorMessage string `json:"errorMessage"`
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	e.ErrorInstance = errors.New(aux.ErrorMessage)
	return nil
}

// Error function for ValidationError to return the error message
func (e ValidationError) Error() string {
	return e.ErrorMessage
}

// NewCommonErrorSomethingWentWrong creates a new CommonError instance
// with `Something went wrong`  ErrorInstance field.
func NewCommonErrorSomethingWentWrong(err error) CommonError {
	return CommonError{
		StatusCode:    http.StatusInternalServerError,
		ErrorCode:     ErrCodeSomethingWentWrong,
		ErrorInstance: err,
	}
}

// NewCommonErrorBadRequest creates a new CommonError instance
// with `Bad request`  ErrorInstance field.
func NewCommonErrorBadRequest(err error) CommonError {
	return CommonError{
		StatusCode:    http.StatusBadRequest,
		ErrorCode:     ErrCodeBadRequest,
		ErrorInstance: err,
	}
}

// IsCommonError returns true if the given error is a CommonError.
func IsCommonError(err error) bool {
	_, ok := err.(CommonError)
	return ok
}

// ParseCommonError checks if the given error is a CommonError instance and returns it.
// If not, it creates a new CommonError instance using NewCommonErrorSomethingWentWrong.
func ParseCommonError(err error) CommonError {
	if e, ok := err.(CommonError); ok {
		return e
	}
	return NewCommonErrorSomethingWentWrong(err)
}

// IsValidationError returns true if the given error is a IsValidationError.
func IsValidationError(err error) bool {
	_, ok := err.(ValidationError)
	return ok
}

// ParseValidationError converts an error into a ValidationError.
// If the input error is a ValidationError, it's returned as is.
func ParseValidationError(err error) ValidationError {
	switch err := err.(type) {
	case ValidationError:
		return err
	case validator.ValidationErrors:
		errDetailList := make([]ValidationErrorDetail, len(err))
		fieldList := make([]string, len(err))
		for i, e := range err {
			fieldList[i] = fmt.Sprintf("'%s'", e.Field())
			errDetailList[i] = ValidationErrorDetail{
				Field:   e.Field(),
				Tag:     e.Tag(),
				Message: fmt.Sprintf(errMessageValidationFailed, e.StructNamespace(), e.Field(), e.Tag()),
			}
		}
		return ValidationError{
			ErrorMessage: fmt.Sprintf("Validation failed for %s", strings.Join(fieldList, ", ")),
			Details:      errDetailList,
			Errors:       err,
		}
	default:
		return ValidationError{
			ErrorMessage: err.Error(),
		}
	}
}

// ParseErrorResponse converts an error into a ErrorResponse.
// If the input error is a ErrorResponse, it's returned as is.
func ParseErrorResponse(err error) ErrorResponse {
	resp := ErrorResponse{
		StatusCode:   http.StatusInternalServerError,
		ErrorCode:    ErrCodeSomethingWentWrong,
		ErrorMessage: err.Error(),
	}
	switch err := err.(type) {
	case *echo.HTTPError:
		resp.StatusCode = err.Code
		switch resp.StatusCode {
		case http.StatusBadRequest:
			resp.ErrorCode = ErrCodeBadRequest
			resp.ErrorMessage = ErrMessageBadRequest
		case http.StatusUnauthorized:
			resp.ErrorCode = ErrCodeUnauthorized
			resp.ErrorMessage = ErrMessageUnauthorized
		case http.StatusNotFound:
			resp.ErrorCode = ErrCodeNotFound
			resp.ErrorMessage = ErrMessageNotFound
		default:
			message := fmt.Sprintf("%v", err.Message)
			if err.Internal != nil {
				message = fmt.Sprintf("%s, %s", message, err.Internal)
			}
			resp.ErrorMessage = message
		}
	case CommonError:
		resp.StatusCode = err.StatusCode
		resp.ErrorCode = err.ErrorCode
		if err.ErrorInstance != nil {
			resp.ErrorMessage = err.ErrorInstance.Error()
			switch e := err.ErrorInstance.(type) {
			case ValidationError, validator.ValidationErrors:
				errValidator := ParseValidationError(e)
				resp.ErrorMessage = errValidator.ErrorMessage
				resp.ErrorValidation = errValidator.Details
			}
		}
	case ValidationError, validator.ValidationErrors:
		errValidator := ParseValidationError(err)
		resp.StatusCode = http.StatusBadRequest
		resp.ErrorCode = ErrCodeBadRequest
		resp.ErrorMessage = errValidator.ErrorMessage
		resp.ErrorValidation = errValidator.Details
		resp.ErrorValidation = errValidator.Details
	}
	return resp
}
