package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		err := Validate.Struct(data)
		assert.Equal(t, test.Expected, err == nil, test.Input)
	}
}

func BenchmarkValidateNumberString(b *testing.B) {
	data := Data{Balance: "10,000.00"}
	for n := 0; n < b.N; n++ {
		Validate.Struct(data)
	}
}
