package utils_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/dollarsignteam/go-utils"
)

func TestIsImageContentType(t *testing.T) {
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
			result := utils.Image.IsImageContentType(fh)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestIsImage(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		content  []byte
		expected bool
	}{
		{
			name:     "valid jpg",
			filename: "test.jpg",
			content:  []byte("\xFF\xD8\xFF"),
			expected: true,
		},
		{
			name:     "valid png",
			filename: "test.png",
			content:  []byte("\x89PNG\x0D\x0A\x1A\x0A"),
			expected: true,
		},
		{
			name:     "invalid file header",
			filename: "test.txt",
			content:  []byte{0x00, 0x01, 0x02, 0x03},
			expected: false,
		},
		{
			name:     "empty file header",
			filename: "test.txt",
			content:  []byte{},
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile("file", test.filename)
			part.Write(test.content)
			writer.Close()
			req := httptest.NewRequest(http.MethodPost, "/", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			c := echo.New().NewContext(req, nil)
			fh, err := c.FormFile("file")
			if assert.NoError(t, err) {
				result := utils.Image.IsImage(fh)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestIsImage_InvalidFile(t *testing.T) {
	fh := &multipart.FileHeader{
		Filename: "invalid-file.txt",
	}
	result := utils.Image.IsImage(fh)
	assert.False(t, result)
	assert.False(t, utils.Image.IsImage(nil))
	assert.False(t, utils.Image.IsImageContentType(nil))
}
