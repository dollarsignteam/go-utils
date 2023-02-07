package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dollarsignteam/go-utils"
)

type Result struct {
	Data string `json:"data" validate:"required"`
}

func TestParseJSON_Error(t *testing.T) {
	result := Result{}
	err := utils.JSON.ParseAndValidate("", &result)
	assert.EqualError(t, err, "unexpected end of JSON input")
}

func TestParseJSON_Pass(t *testing.T) {
	result := Result{}
	err := utils.JSON.ParseAndValidate(`{"data": "ok"}`, &result)
	assert.Nil(t, err)
	assert.Equal(t, "ok", result.Data)
}

func TestParseJSON_Failed(t *testing.T) {
	result := Result{}
	err := utils.JSON.ParseAndValidate("{}", &result)
	assert.EqualError(t, err, "Key: 'Result.Data' Error:Field validation for 'Data' failed on the 'required' tag")
}
