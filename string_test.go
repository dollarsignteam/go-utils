package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dollarsignteam/go-utils"
)

func TestRemoveDuplicateSpaces(t *testing.T) {
	s := "  Hello \t\r\n - World  \u200B!  "
	expected := "Hello - World !"
	result := utils.String.RemoveDuplicateSpaces(s)
	assert.Equal(t, expected, result)
}

func TestRemoveAllSpaces(t *testing.T) {
	s := "  Hello \t\r\n - World  \u200B!  "
	expected := "Hello-World!"
	result := utils.String.RemoveAllSpaces(s)
	assert.Equal(t, expected, result)
}
