package utils

import "golang.org/x/exp/constraints"

// Map applies a function to each element of a slice and returns a new slice
//
//	s is the slice to be mapped
//	f is the function to be applied to each element of the slice
func Map[T any](s []T, f func(T) T) []T {
	result := make([]T, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

// Reduce applies a function to each element of a slice and return the accumulated result
//
//	s is the slice to be reduced
//	f is the function to be applied to each element of the slice
//	initial is the initial value of the accumulator
func Reduce[T any](s []T, f func(prev T, curr T) T, initial T) T {
	r := initial
	for _, v := range s {
		r = f(r, v)
	}
	return r
}

// Filter applies a function to each element of a slice and returns a new slice with the elements that passed the test
//
//	s is the slice to be filtered
//	f is the function to be applied to each element of the slice
func Filter[T any](s []T, f func(T) bool) []T {
	r := make([]T, 0)
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

// SimpleSort sorts a slice of ordered elements using efficient merge sort
//
//	Note: this function is not generic, it only works with ordered elements
//	s is the slice to be sorted
//	returns a new slice with the elements sorted
func SimpleSort[T constraints.Ordered](s []T) []T {
	result := make([]T, len(s))
	copy(result, s)
	if len(result) < 2 {
		return result
	}
	mid := len(result) / 2
	return merge(SimpleSort(result[:mid]), SimpleSort(result[mid:]))
}

// merge is a helper function for SimpleSort
// it merges two slices of ordered elements
// it is not generic, it only works with ordered elements
//
//	a and b are the slices to be merged
//	returns a new slice with the elements merged
func merge[T constraints.Ordered](a, b []T) []T {
	result := make([]T, len(a)+len(b))
	i := 0
	for len(a) > 0 && len(b) > 0 {
		if a[0] < b[0] {
			result[i] = a[0]
			a = a[1:]
		} else {
			result[i] = b[0]
			b = b[1:]
		}
		i++
	}
	for j := 0; j < len(a); j++ {
		result[i] = a[j]
		i++
	}
	for j := 0; j < len(b); j++ {
		result[i] = b[j]
		i++
	}
	return result
}

// SortFunc sorts a slice of elements using a function to compare them
func SortFunc[T any](s []T, f func(T, T) bool) []T {
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if f(s[i], s[j]) {
				s[i], s[j] = s[j], s[i]
			}
		}
	}
	return s
}
