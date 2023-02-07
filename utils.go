package utils

import (
	"runtime"
	"strings"
)

func PointerOf[T any](v T) *T {
	return &v
}

func PackageName() string {
	pc, _, _, _ := runtime.Caller(1)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	pkg := strings.Split(parts[len(parts)-1], ".")
	return pkg[0]
}
