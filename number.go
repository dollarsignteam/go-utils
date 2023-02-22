package utils

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

// Min returns the minimum of two ordered values x and y.
func Min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

// Max returns the maximum of two ordered values x and y.
func Max[T constraints.Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
}

// MinOf returns the minimum value in the ordered slice s.
// Returns the zero value of T if s is empty.
func MinOf[T constraints.Ordered](s []T) T {
	if len(s) == 0 {
		var zero T
		return zero
	}
	m := s[0]
	for _, v := range s {
		if m > v {
			m = v
		}
	}
	return m
}

// MaxOf returns the maximum value in the ordered slice s.
// Returns the zero value of T if s is empty.
func MaxOf[T constraints.Ordered](s []T) T {
	if len(s) == 0 {
		var zero T
		return zero
	}
	m := s[0]
	for _, v := range s {
		if m < v {
			m = v
		}
	}
	return m
}

// RandomInt64 generates a random int64 value between min and max.
func RandomInt64(min, max int64) int64 {
	if min > max {
		min, max = max, min
	}
	nBig, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
	return nBig.Int64() + min
}

// RandomFloat64 generates a random float64 value between min and max.
func RandomFloat64(min, max float64) float64 {
	if min > max {
		min, max = max, min
	}
	nBig, _ := rand.Int(rand.Reader, big.NewInt(1<<62))
	return (float64(nBig.Int64())/float64(1<<62))*(max-min) + min
}

// ParseFloat64 parses a string s as a float64 value.
func ParseFloat64(s string) (float64, error) {
	s = strings.ReplaceAll(strings.TrimSpace(s), ",", "")
	return strconv.ParseFloat(s, 64)
}
