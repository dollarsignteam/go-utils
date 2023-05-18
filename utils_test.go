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
		input    []any
		expected []any
	}{
		{
			input:    []any{1, 2, 3, 2, 1},
			expected: []any{1, 2, 3},
		},
		{
			input:    []any{"foo", "bar", "baz", "foo"},
			expected: []any{"foo", "bar", "baz"},
		},
		{
			input:    []any{1, "foo", true, 2, "foo"},
			expected: []any{1, "foo", true, 2},
		},
	}
	for _, test := range tests {
		result := utils.UniqueOf(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestValueOf(t *testing.T) {
	var i *int
	result := utils.ValueOf(i)
	assert.Equal(t, 0, result)

	value := 5
	result = utils.ValueOf(&value)
	assert.Equal(t, value, result)

	var s *string
	resultStr := utils.ValueOf(s)
	assert.Equal(t, "", resultStr)

	str := "hello"
	resultStr = utils.ValueOf(&str)
	assert.Equal(t, str, resultStr)

	var f *float64
	resultFloat := utils.ValueOf(f)
	assert.Equal(t, 0.0, resultFloat)

	floatVal := 3.14
	resultFloat = utils.ValueOf(&floatVal)
	assert.Equal(t, floatVal, resultFloat)
}

func TestIsArrayOrSlice(t *testing.T) {
	slice := []int{1, 2, 3}
	array := [3]int{4, 5, 6}
	tests := []struct {
		Input    any
		Expected bool
	}{
		{Input: nil, Expected: false},
		{Input: slice, Expected: true},
		{Input: array, Expected: true},
		{Input: &slice, Expected: true},
		{Input: &array, Expected: true},
	}
	for _, test := range tests {
		result := utils.IsArrayOrSlice(test.Input)
		assert.Equal(t, test.Expected, result)
	}
}

func TestBoolToInt(t *testing.T) {
	tests := []struct {
		input    bool
		expected int
	}{
		{input: true, expected: 1},
		{input: false, expected: 0},
	}
	for _, test := range tests {
		result := utils.BoolToInt(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestIntToBool(t *testing.T) {
	tests := []struct {
		input    int
		expected bool
	}{
		{input: 1, expected: true},
		{input: 0, expected: false},
		{input: -1, expected: true},
	}
	for _, test := range tests {
		result := utils.IntToBool(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestUnionOf(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]any
		expected []any
	}{
		{
			name:     "no input",
			input:    [][]any{},
			expected: []any{},
		},
		{
			name:     "one empty slice",
			input:    [][]any{{}},
			expected: []any{},
		},
		{
			name:     "one non-empty slice",
			input:    [][]any{{1, 2, 3}},
			expected: []any{1, 2, 3},
		},
		{
			name:     "two non-empty slices",
			input:    [][]any{{1, 2, 3, 3, true}, {3, 4, "foo"}},
			expected: []any{1, 2, 3, true, 4, "foo"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.UnionOf(test.input...)
			assert.ElementsMatch(t, test.expected, result)
		})
	}
}

func TestIntersectionOf(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]int
		expected []int
	}{
		{
			name:     "empty input",
			input:    nil,
			expected: []int{},
		},
		{
			name:     "single slice",
			input:    [][]int{{1, 2, 3}},
			expected: []int{1, 2, 3},
		},
		{
			name:     "multiple slices with one intersection",
			input:    [][]int{{1, 2}, {2, 3}, {2, 4}},
			expected: []int{2},
		},
		{
			name:     "multiple slices with no intersection",
			input:    [][]int{{1, 2}, {2, 3}, {4, 5}},
			expected: []int{},
		},
		{
			name:     "multiple slices with multiple intersections",
			input:    [][]int{{1, 2, 3}, {2, 3}, {3, 4}},
			expected: []int{3},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.IntersectionOf(test.input...)
			assert.ElementsMatch(t, test.expected, result)
		})
	}
}

func TestDifferenceOf(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]int
		expected []int
	}{
		{
			name:     "empty input",
			input:    [][]int{},
			expected: []int{},
		},
		{
			name:     "single slice input",
			input:    [][]int{{1, 2, 3}},
			expected: []int{1, 2, 3},
		},
		{
			name:     "two slice inputs with no common elements",
			input:    [][]int{{1, 2, 3}, {4, 5, 6}},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:     "two slice inputs with some common elements",
			input:    [][]int{{1, 2, 3}, {3, 4, 5}},
			expected: []int{1, 2, 4, 5},
		},
		{
			name:     "three slice inputs with some common elements",
			input:    [][]int{{1, 2, 3}, {2, 3, 4}, {3, 4, 5}},
			expected: []int{1, 5},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.DifferenceOf(test.input...)
			assert.ElementsMatch(t, test.expected, result)
		})
	}
}

func BenchmarkUniqueOf(b *testing.B) {
	for n := 0; n < b.N; n++ {
		utils.UniqueOf([]any{1, "foo", true, 2, "foo"})
	}
}

func BenchmarkUnionOf(b *testing.B) {
	for n := 0; n < b.N; n++ {
		utils.UnionOf([][]any{{1, 2, 3, 3, true}, {3, 4, "foo"}}...)
	}
}

func BenchmarkIntersectionOf(b *testing.B) {
	for n := 0; n < b.N; n++ {
		utils.IntersectionOf([][]int{{1, 2, 3}, {2, 3}, {3, 4}}...)
	}
}

func BenchmarkDifferenceOf(b *testing.B) {
	for n := 0; n < b.N; n++ {
		utils.DifferenceOf([][]int{{1, 2, 3}, {2, 3, 4}, {3, 4, 5}}...)
	}
}

func BenchmarkIsArrayOrSlice(b *testing.B) {
	slice := []int{1, 2, 3}
	for n := 0; n < b.N; n++ {
		utils.IsArrayOrSlice(&slice)
	}
}
