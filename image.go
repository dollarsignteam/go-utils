package utils

import (
	"mime/multipart"
	"strings"
)

// Image utility instance
var Image imageUtil

// imageUtil is a utility struct for image related functions
type imageUtil struct{}

// IsImage checks if the file is an image
func (i imageUtil) IsImage(fh *multipart.FileHeader) bool {
	contentType := fh.Header.Get("Content-Type")
	return strings.HasPrefix(contentType, "image/")
}
