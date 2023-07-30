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

func TestHashAndVerifyPassword(t *testing.T) {
	password := "foo"
	result, err := utils.String.HashPassword(password)
	assert.Nil(t, err)
	err = utils.String.VerifyPassword(result, password)
	assert.Nil(t, err)
}

func TestHashPassword_Error(t *testing.T) {
	password := make([]byte, 80)
	result, err := utils.String.HashPassword(string(password))
	assert.Empty(t, result)
	assert.EqualError(t, err, "bcrypt: password length exceeds 72 bytes")
}

func TestHashCrc32(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "Test Case 1", input: "Hello, World!", expected: "ec4ac3d0"},
		{name: "Test Case 2", input: "Lorem Ipsum", expected: "358ad45d"},
		{name: "Test Case 3", input: "1234567890", expected: "261daee5"},
		{name: "Test Case 4", input: "A4D7B2B7-D62D-423C-B0C2-2A871F98E427", expected: "042b1405"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.String.HashCrc32(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
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

func BenchmarkHashPassword(b *testing.B) {
	for n := 0; n < b.N; n++ {
		utils.String.HashPassword(TestSpacesString)
	}
}

func BenchmarkHashCrc32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utils.String.HashCrc32(TestSpacesString)
	}
}
