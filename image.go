package utils

import (
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

// Image utility instance
var Image ImageUtil

// ImageUtil is a utility struct for image related functions
type ImageUtil struct{}

// IsImageContentType checks if the content type is an image
func (ImageUtil) IsImageContentType(fh *multipart.FileHeader) bool {
	if fh == nil {
		return false
	}
	contentType := fh.Header.Get("Content-Type")
	return strings.HasPrefix(contentType, "image/")
}

// IsImage checks if the file is an image
func (ImageUtil) IsImage(fh *multipart.FileHeader) bool {
	if fh == nil {
		return false
	}
	f, err := fh.Open()
	if err != nil {
		return false
	}
	defer f.Close()
	buf := make([]byte, 512)
	if _, err = io.ReadFull(f, buf); err != nil && err != io.ErrUnexpectedEOF {
		return false
	}
	contentType := http.DetectContentType(buf)
	return strings.HasPrefix(contentType, "image/")
}
