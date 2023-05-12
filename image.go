package utils

import (
	"mime/multipart"
	"strings"
)

// Image utility instance
var Image ImageUtil

// ImageUtil is a utility struct for image related functions
type ImageUtil struct{}

// IsImage checks if the file is an image
func (ImageUtil) IsImage(fh *multipart.FileHeader) bool {
	contentType := fh.Header.Get("Content-Type")
	return strings.HasPrefix(contentType, "image/")
}
