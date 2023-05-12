package utils_test

import (
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dollarsignteam/go-utils"
)

func TestIsImage(t *testing.T) {
	tests := []struct {
		name     string
		fileType string
		expected bool
	}{
		{"valid image", "image/jpeg", true},
		{"invalid image", "text/plain", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fh := &multipart.FileHeader{
				Header: map[string][]string{
					"Content-Type": {test.fileType},
				},
			}
			result := utils.Image.IsImage(fh)
			assert.Equal(t, test.expected, result)
		})
	}
}
