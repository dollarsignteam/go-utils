package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dollarsignteam/go-utils"
)

func TestMin(t *testing.T) {
	testCases := []struct {
		x        float64
		y        float64
		expected float64
	}{
		{x: 0, y: 1, expected: 0},
		{x: 10, y: 1, expected: 1},
	}
	for _, tc := range testCases {
		min := utils.Min(tc.x, tc.y)
		assert.Equal(t, tc.expected, min)
	}
}

func TestMax(t *testing.T) {
	testCases := []struct {
		x        float64
		y        float64
		expected float64
	}{
		{x: 0, y: 1, expected: 1},
		{x: 10, y: 1, expected: 10},
	}
	for _, tc := range testCases {
		min := utils.Max(tc.x, tc.y)
		assert.Equal(t, tc.expected, min)
	}
}

func TestMinOf(t *testing.T) {
	testCases := []struct {
		s        []float64
		expected float64
	}{
		{s: []float64{}, expected: 0},
		{s: []float64{10, 2, 4, 1, 6, 8, 2}, expected: 1},
	}
	for _, tc := range testCases {
		min := utils.MinOf(tc.s)
		assert.Equal(t, tc.expected, min)
	}
}

func TestMaxOf(t *testing.T) {
	testCases := []struct {
		s        []float64
		expected float64
	}{
		{s: []float64{}, expected: 0},
		{s: []float64{1, 2, 4, 10, 6, 8, 2}, expected: 10},
	}
	for _, tc := range testCases {
		min := utils.MaxOf(tc.s)
		assert.Equal(t, tc.expected, min)
	}
}
