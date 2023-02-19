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

func UniqueOf[T any](input []T) []T {
	u := make([]T, 0, len(input))
	m := make(map[any]struct{})
	for _, v := range input {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			u = append(u, v)
		}
	}
	return u
}

func ValueOf[T any](ptr *T) T {
	if ptr == nil {
		var v T
		return v
	}
	return *ptr
}
