package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dollarsignteam/go-utils"
)

func TestMin(t *testing.T) {
	tests := []struct {
		x        float64
		y        float64
		expected float64
	}{
		{x: 0, y: 1, expected: 0},
		{x: 10, y: 1, expected: 1},
	}
	for _, test := range tests {
		min := utils.Min(test.x, test.y)
		assert.Equal(t, test.expected, min)
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		x        float64
		y        float64
		expected float64
	}{
		{x: 0, y: 1, expected: 1},
		{x: 10, y: 1, expected: 10},
	}
	for _, test := range tests {
		min := utils.Max(test.x, test.y)
		assert.Equal(t, test.expected, min)
	}
}

func TestMinOf(t *testing.T) {
	tests := []struct {
		input    []float64
		expected float64
	}{
		{input: []float64{}, expected: 0},
		{input: []float64{10, 2, 4, 1, 6, 8, 2}, expected: 1},
	}
	for _, test := range tests {
		min := utils.MinOf(test.input)
		assert.Equal(t, test.expected, min)
	}
}

func TestMaxOf(t *testing.T) {
	tests := []struct {
		input    []float64
		expected float64
	}{
		{input: []float64{}, expected: 0},
		{input: []float64{1, 2, 4, 10, 6, 8, 2}, expected: 10},
	}
	for _, test := range tests {
		min := utils.MaxOf(test.input)
		assert.Equal(t, test.expected, min)
	}
}

func TestRandomInt64(t *testing.T) {
	tests := []struct {
		min         int64
		max         int64
		expectedMin int64
		expectedMax int64
	}{
		{min: 5, max: 10, expectedMin: 5, expectedMax: 10},
		{min: 5, max: -10, expectedMin: -10, expectedMax: 5},
	}
	for _, test := range tests {
		result := utils.RandomInt64(test.min, test.max)
		assert.GreaterOrEqual(t, result, test.expectedMin)
		assert.LessOrEqual(t, result, test.expectedMax)
	}
}

func TestRandomFloat64(t *testing.T) {
	tests := []struct {
		min         float64
		max         float64
		expectedMin float64
		expectedMax float64
	}{
		{min: 5, max: 10, expectedMin: 5, expectedMax: 10},
		{min: 5, max: -10, expectedMin: -10, expectedMax: 5},
	}
	for _, test := range tests {
		result := utils.RandomFloat64(test.min, test.max)
		assert.GreaterOrEqual(t, result, test.expectedMin)
		assert.LessOrEqual(t, result, test.expectedMax)
	}
}

func TestParseFloat64(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{input: "123.45", expected: 123.45},
		{input: "0.1234", expected: 0.1234},
		{input: "-123.45", expected: -123.45},
		{input: "-0.1234", expected: -0.1234},
		{input: "+123.45", expected: 123.45},
		{input: "+0.1234", expected: 0.1234},
		{input: " 9,999,999,999.99 ", expected: 9999999999.99},
		{input: "message", expected: 0},
	}
	for _, test := range tests {
		result, _ := utils.ParseFloat64(test.input)
		assert.Equal(t, test.expected, result)
	}
}
