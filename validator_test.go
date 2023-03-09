package utils_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dollarsignteam/go-utils"
)

type Data struct {
	Balance string `validate:"required,number_string"`
}

func TestValidateNumberString(t *testing.T) {
	tests := []struct {
		Input    string
		Expected bool
	}{
		{Input: "123", Expected: true},
		{Input: "-123", Expected: true},
		{Input: "1.23", Expected: true},
		{Input: "-1.23", Expected: true},
		{Input: "1,000.00", Expected: true},
		{Input: "-1,000.99", Expected: true},
		{Input: ".123456789", Expected: true},
		{Input: "-1,0,00.99", Expected: false},
		{Input: "-10,00.99", Expected: false},
		{Input: "1.x", Expected: false},
	}
	for _, test := range tests {
		data := Data{Balance: test.Input}
		err := utils.Validate.Struct(data)
		assert.Equal(t, test.Expected, err == nil, test.Input)
	}
}

func TestGetJSONTagName(t *testing.T) {
	type User struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"-"`
	}
	tests := []struct {
		fieldName string
		expected  string
	}{
		{fieldName: "ID", expected: "id"},
		{fieldName: "Name", expected: "name"},
		{fieldName: "Email", expected: ""},
	}
	for _, test := range tests {
		field, _ := reflect.TypeOf(User{}).FieldByName(test.fieldName)
		result := utils.GetJSONTagName(field)
		assert.Equal(t, test.expected, result)
	}
}

func TestValidateStruct_Valid(t *testing.T) {
	s := Data{
		Balance: "123",
	}
	err := utils.ValidateStruct(s)
	assert.NoError(t, err)
}

func TestValidateStruct_Invalid(t *testing.T) {
	s := Data{}
	err := utils.ValidateStruct(s)
	assert.EqualError(t, err, "Validation failed for 'Balance'")
}

func BenchmarkValidateNumberString(b *testing.B) {
	data := Data{Balance: "10,000.00"}
	for n := 0; n < b.N; n++ {
		utils.Validate.Struct(data)
	}
}
