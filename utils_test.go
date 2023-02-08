package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dollarsignteam/go-utils"
)

func TestPointerOf(t *testing.T) {
	value := 0
	result := utils.PointerOf(value)
	assert.Equal(t, &value, result)
}

func TestPackageName(t *testing.T) {
	name := utils.PackageName()
	assert.Equal(t, "go-utils_test", name)
}

func TestUniqueOf(t *testing.T) {
	tests := []struct {
		input    []interface{}
		expected []interface{}
	}{
		{
			input:    []interface{}{1, 2, 3, 2, 1},
			expected: []interface{}{1, 2, 3},
		},
		{
			input:    []interface{}{"foo", "bar", "baz", "foo"},
			expected: []interface{}{"foo", "bar", "baz"},
		},
		{
			input:    []interface{}{1, "foo", true, 2, "foo"},
			expected: []interface{}{1, "foo", true, 2},
		},
	}
	for _, test := range tests {
		result := utils.UniqueOf(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func BenchmarkUniqueOf(b *testing.B) {
	for n := 0; n < b.N; n++ {
		utils.UniqueOf([]interface{}{1, "foo", true, 2, "foo"})
	}
}
