package utils

import (
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

var (
	ErrNotAnImage = errors.New("not an image")
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

func (i ImageUtil) FromMultipart(fh *multipart.FileHeader) (image.Image, error) {
	if !i.IsImage(fh) {
		return nil, ErrNotAnImage
	}
	file, err := fh.Open()
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(file)
	return img, err
}

func (i ImageUtil) QRStringReader(img image.Image) (string, error) {
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}
