package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dollarsignteam/go-utils"
)

var TestSpacesString = "  foo \t\r\n - bar  \u200B!  "

func TestRemoveDuplicateSpaces(t *testing.T) {
	expected := "foo - bar !"
	result := utils.String.RemoveDuplicateSpaces(TestSpacesString)
	assert.Equal(t, expected, result)
}

func TestRemoveAllSpaces(t *testing.T) {
	expected := "foo-bar!"
	result := utils.String.RemoveAllSpaces(TestSpacesString)
	assert.Equal(t, expected, result)
}

func TestUUID(t *testing.T) {
	result := utils.String.UUID()
	assert.Len(t, result, 36)
}

func TestMD5(t *testing.T) {
	result := utils.String.MD5(TestSpacesString)
	expected := "34130b8b17f2e67b2da09cd24f868885"
	assert.Equal(t, expected, result)
}

func TestSHA1(t *testing.T) {
	result := utils.String.SHA1(TestSpacesString)
	expected := "5045a76071cb10b5fa05b73af1e8e174d4979955"
	assert.Equal(t, expected, result)
}

func TestSHA256(t *testing.T) {
	result := utils.String.SHA256(TestSpacesString)
	expected := "bafa334ba4639eca91f087ad98a0dcc9d1ac2f82da8beafed8fbaad717a51c6d"
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

func BenchmarkMD5(b *testing.B) {
	for n := 0; n < b.N; n++ {
		utils.String.MD5(TestSpacesString)
	}
}

func BenchmarkSHA1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		utils.String.SHA1(TestSpacesString)
	}
}

func BenchmarkSHA256(b *testing.B) {
	for n := 0; n < b.N; n++ {
		utils.String.SHA256(TestSpacesString)
	}
}
