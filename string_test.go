package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dollarsignteam/go-utils"
)

var TestSpacesString = "  Hello \t\r\n - World  \u200B!  "

func TestRemoveDuplicateSpaces(t *testing.T) {
	expected := "Hello - World !"
	result := utils.String.RemoveDuplicateSpaces(TestSpacesString)
	assert.Equal(t, expected, result)
}

func TestRemoveAllSpaces(t *testing.T) {
	expected := "Hello-World!"
	result := utils.String.RemoveAllSpaces(TestSpacesString)
	assert.Equal(t, expected, result)
}

func BenchmarkRemoveDuplicateSpaces(b *testing.B) {
	for n := 0; n < b.N; n++ {
		utils.String.RemoveDuplicateSpaces(TestSpacesString)
	}
}

func BenchmarkRemoveAllSpaces(b *testing.B) {
	for n := 0; n < b.N; n++ {
		utils.String.RemoveAllSpaces(TestSpacesString)
	}
}
