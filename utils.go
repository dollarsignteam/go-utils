package utils

import (
	"reflect"
	"runtime"
	"strings"

	"golang.org/x/exp/slices"
)

// PointerOf returns a pointer to the input value
func PointerOf[T any](v T) *T {
	return &v
}

// PackageName returns the name of the package that calls it.
func PackageName() string {
	pc, _, _, _ := runtime.Caller(1)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	pkg := strings.Split(parts[len(parts)-1], ".")
	return pkg[0]
}

// UniqueOf removes duplicates from a slice of any type
// and returns a new slice containing only the unique elements.
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

// ValueOf takes a pointer to a value of any type and returns the value.
func ValueOf[T any](ptr *T) T {
	if ptr == nil {
		var v T
		return v
	}
	return *ptr
}

// IsArrayOrSlice takes a value of any type
// and returns a boolean indicating if it is a slice or an array.
func IsArrayOrSlice(i any) bool {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return (v.Kind() == reflect.Slice || v.Kind() == reflect.Array)
}

// BoolToInt converts a boolean value to an integer
// (1 for true, 0 for false).
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// IntToBool converts an integer value to a boolean
// (true for non-zero values, false for zero).
func IntToBool(i int) bool {
	return i != 0
}

// UnionOf returns the union of all input slices,
// preserving the order of the first occurrence of each unique element.
func UnionOf[T any](input ...[]T) []T {
	m := make(map[any]struct{})
	for _, s := range input {
		for _, v := range s {
			m[v] = struct{}{}
		}
	}
	u := make([]T, 0, len(m))
	for k := range m {
		u = append(u, k.(T))
	}
	return u
}

// IntersectionOf returns the intersection of all input slices,
// preserving the order of the first occurrence of each shared element.
func IntersectionOf[T comparable](input ...[]T) []T {
	if len(input) == 0 {
		return []T{}
	}
	if len(input) == 1 {
		return input[0]
	}
	result := []T{}
	for _, v := range input[0] {
		foundInAll := true
		for i := 1; i < len(input); i++ {
			if !slices.Contains(input[i], v) {
				foundInAll = false
				break
			}
		}
		if foundInAll {
			result = append(result, v)
		}
	}
	return result
}

// DifferenceOf returns a slice containing all unique elements across all input slices,
// The order of elements in the resulting slice is not guaranteed
func DifferenceOf[T comparable](input ...[]T) []T {
	if len(input) == 0 {
		return []T{}
	}
	if len(input) == 1 {
		return input[0]
	}
	elementCounts := make(map[T]int)
	for _, inputSlice := range input {
		for _, element := range inputSlice {
			elementCounts[element]++
		}
	}
	result := []T{}
	for element, count := range elementCounts {
		if count == 1 {
			result = append(result, element)
		}
	}
	return result
}
