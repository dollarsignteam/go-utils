package utils

import (
	"runtime"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func PointerOf[T any](v T) *T {
	return &v
}

func PackageName() string {
	pc, _, _, _ := runtime.Caller(1)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	pkg := strings.Split(parts[len(parts)-1], ".")
	return pkg[0]
}
