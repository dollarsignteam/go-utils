package utils_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/dollarsignteam/go-utils"
)

func TestCommonError_Error(t *testing.T) {
	errMessage := "Something went wrong"
	commonErr := utils.CommonError{
		StatusCode: 500,
		ErrorCode:  "ERROR",
	}
	result := commonErr.Error()
	assert.Equal(t, errMessage, result)
}

func TestCommonError_MarshalJSON(t *testing.T) {
	errMessage := "Internal error"
	err := errors.New(errMessage)
	commonErr := utils.CommonError{
		StatusCode:    500,
		ErrorCode:     "ERROR",
		ErrorInstance: err,
	}
	result, err := json.Marshal(commonErr)
	assert.NoError(t, err)
	expected := `{"statusCode":500,"errorCode":"ERROR","errorMessage":"Internal error"}`
	assert.JSONEq(t, expected, string(result))
}

func TestCommonError_UnmarshalJSON(t *testing.T) {
	jsonStr := `{"statusCode":500,"errorCode":"ERROR","errorMessage":"Internal error"}`
	var commonErr utils.CommonError
	err := json.Unmarshal([]byte(jsonStr), &commonErr)
	assert.NoError(t, err)
	assert.Equal(t, 500, commonErr.StatusCode)
	assert.Equal(t, "ERROR", commonErr.ErrorCode)
	assert.Equal(t, "Internal error", commonErr.Error())
}

func TestCommonError_UnmarshalJSON_InvalidError(t *testing.T) {
	jsonStr := `[]`
	var commonErr utils.CommonError
	err := json.Unmarshal([]byte(jsonStr), &commonErr)
	assert.Error(t, err)
	assert.Equal(t, "Something went wrong", commonErr.Error())
}

func TestValidationError_Error(t *testing.T) {
	errMessage := "Validation failed"
	validationErr := utils.ValidationError{
		ErrorMessage: errMessage,
	}
	result := validationErr.Error()
	assert.Equal(t, errMessage, result)
}

func TestValidationErrorDetail(t *testing.T) {
	detail := utils.ValidationErrorDetail{
		Field:   "id",
		Tag:     "required",
		Message: "Key: 'Member.id' Error:Field validation for 'id' failed on the 'required' tag'",
	}
	field := detail.Field
	tag := detail.Tag
	message := detail.Message
	assert.Equal(t, "id", field)
	assert.Equal(t, "required", tag)
	assert.Equal(t, "Key: 'Member.id' Error:Field validation for 'id' failed on the 'required' tag'", message)
}

func TestNewCommonErrorSomethingWentWrong(t *testing.T) {
	err := errors.New("test error")
	commonErr := utils.NewCommonErrorSomethingWentWrong(err)
	assert.Equal(t, http.StatusInternalServerError, commonErr.StatusCode)
	assert.Equal(t, utils.ErrCodeSomethingWentWrong, commonErr.ErrorCode)
	assert.Equal(t, err, commonErr.ErrorInstance)
	assert.Equal(t, err.Error(), commonErr.Error())
}

func TestNewCommonErrorBadRequest(t *testing.T) {
	err := errors.New("test error")
	commonErr := utils.NewCommonErrorBadRequest(err)
	assert.Equal(t, http.StatusBadRequest, commonErr.StatusCode)
	assert.Equal(t, utils.ErrCodeBadRequest, commonErr.ErrorCode)
	assert.Equal(t, err, commonErr.ErrorInstance)
	assert.Equal(t, err.Error(), commonErr.Error())
}

func TestParseCommonError(t *testing.T) {
	err := errors.New("test error")
	commonErr := utils.CommonError{
		StatusCode:    http.StatusBadRequest,
		ErrorCode:     "BadRequest",
		ErrorInstance: err,
	}
	assert.Equal(t, commonErr, utils.ParseCommonError(commonErr))
	newCommonErr := utils.ParseCommonError(err)
	assert.Equal(t, http.StatusInternalServerError, newCommonErr.StatusCode)
	assert.Equal(t, utils.ErrCodeSomethingWentWrong, newCommonErr.ErrorCode)
	assert.Equal(t, err, newCommonErr.ErrorInstance)
	assert.Equal(t, err.Error(), newCommonErr.Error())
}

func TestParseValidationError(t *testing.T) {
	t.Run("ValidationError", func(t *testing.T) {
		err := utils.ValidationError{
			ErrorMessage: "test error message",
		}
		result := utils.ParseValidationError(err)
		assert.Equal(t, err, result)
	})

	t.Run("ValidationErrors", func(t *testing.T) {
		data := Data{Balance: "foo"}
		err := utils.Validate.Struct(data)
		result := utils.ParseValidationError(err)
		expected := utils.ValidationError{
			ErrorMessage: "Validation failed for 'Balance'",
			Details: []utils.ValidationErrorDetail{
				{
					Field:   "Balance",
					Tag:     "number_string",
					Message: "Key: 'Data.Balance', Error: Validation for 'Balance' failed on the 'number_string' tag",
				},
			},
			Errors: err.(validator.ValidationErrors),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("UnknownError", func(t *testing.T) {
		errMessage := "test error message"
		err := errors.New(errMessage)
		result := utils.ParseValidationError(err)
		assert.EqualError(t, result, errMessage)
	})
}

func TestIsCommonError(t *testing.T) {
	tests := []struct {
		Input    error
		Expected bool
	}{
		{
			Input: utils.CommonError{
				StatusCode: 500,
				ErrorCode:  "ERROR",
			},
			Expected: true,
		},
		{
			Input:    errors.New("regular error"),
			Expected: false,
		},
	}

	for _, test := range tests {
		result := utils.IsCommonError(test.Input)
		assert.Equal(t, test.Expected, result)
	}
}

type UnknownError struct {
	utils.CommonError
}

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		Input    error
		Expected bool
	}{
		{
			Input: utils.ValidationError{
				ErrorMessage: "validation error",
			},
			Expected: true,
		},
		{
			Input:    errors.New("regular error"),
			Expected: false,
		},
		{
			Input:    UnknownError{},
			Expected: false,
		},
	}

	for _, test := range tests {
		result := utils.IsValidationError(test.Input)
		assert.Equal(t, test.Expected, result)
	}
}

func TestParseErrorResponse(t *testing.T) {
	t.Run("HTTPErrorBadRequest", func(t *testing.T) {
		err := echo.ErrBadRequest
		resp := utils.ParseErrorResponse(err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, utils.ErrCodeBadRequest, resp.ErrorCode)
		assert.Equal(t, utils.ErrMessageBadRequest, resp.ErrorMessage)
	})

	t.Run("HTTPErrorUnauthorized", func(t *testing.T) {
		err := echo.ErrUnauthorized
		resp := utils.ParseErrorResponse(err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		assert.Equal(t, utils.ErrCodeUnauthorized, resp.ErrorCode)
		assert.Equal(t, utils.ErrMessageUnauthorized, resp.ErrorMessage)
	})

	t.Run("HTTPErrorNotFound", func(t *testing.T) {
		err := echo.ErrNotFound
		resp := utils.ParseErrorResponse(err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.Equal(t, utils.ErrCodeNotFound, resp.ErrorCode)
		assert.Equal(t, utils.ErrMessageNotFound, resp.ErrorMessage)
	})

	t.Run("HTTPErrorForbidden", func(t *testing.T) {
		err := echo.ErrForbidden
		resp := utils.ParseErrorResponse(err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		assert.Equal(t, utils.ErrCodeForbidden, resp.ErrorCode)
		assert.Equal(t, utils.ErrMessageForbidden, resp.ErrorMessage)
	})

	t.Run("HTTPErrorMethodNotAllowed", func(t *testing.T) {
		err := echo.ErrMethodNotAllowed
		resp := utils.ParseErrorResponse(err)
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
		assert.Equal(t, utils.ErrCodeMethodNotAllowed, resp.ErrorCode)
		assert.Equal(t, utils.ErrMessageMethodNotAllowed, resp.ErrorMessage)
	})

	t.Run("HTTPErrorInternalError", func(t *testing.T) {
		err := &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "Oops, something went wrong!",
			Internal: errors.New("internal error"),
		}
		expected := "Oops, something went wrong!, internal error"
		resp := utils.ParseErrorResponse(err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, utils.ErrCodeSomethingWentWrong, resp.ErrorCode)
		assert.Equal(t, expected, resp.ErrorMessage)
	})

	t.Run("CommonErrorWithValidationError", func(t *testing.T) {
		data := Data{Balance: "foo"}
		errValidation := utils.ValidateStruct(data)
		errMessage := "Validation failed for 'Balance'"
		err := utils.CommonError{
			StatusCode:    http.StatusBadRequest,
			ErrorCode:     utils.ErrCodeBadRequest,
			ErrorInstance: errValidation,
		}
		resp := utils.ParseErrorResponse(err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, utils.ErrCodeBadRequest, resp.ErrorCode)
		assert.Equal(t, errMessage, resp.ErrorMessage)
		assert.Equal(t, errValidation.(utils.ValidationError).Details, resp.ErrorValidation)
	})

	t.Run("CommonErrorWithValidationErrors", func(t *testing.T) {
		data := Data{Balance: "foo"}
		errValidation := utils.Validate.Struct(data)
		errMessage := "Validation failed for 'Balance'"
		err := utils.CommonError{
			StatusCode:    http.StatusBadRequest,
			ErrorCode:     utils.ErrCodeBadRequest,
			ErrorInstance: errValidation,
		}
		resp := utils.ParseErrorResponse(err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, utils.ErrCodeBadRequest, resp.ErrorCode)
		assert.Equal(t, errMessage, resp.ErrorMessage)
	})

	t.Run("ValidationErrorAndValidationErrors", func(t *testing.T) {
		data := Data{Balance: "foo"}
		errMessage := "Validation failed for 'Balance'"
		errValidation := utils.Validate.Struct(data)
		resp := utils.ParseErrorResponse(errValidation)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, utils.ErrCodeBadRequest, resp.ErrorCode)
		assert.Equal(t, errMessage, resp.ErrorMessage)
	})
}

func BenchmarkParseValidationError(b *testing.B) {
	data := Data{Balance: "foo"}
	err := utils.Validate.Struct(data)
	for n := 0; n < b.N; n++ {
		utils.ParseValidationError(err)
	}
}
