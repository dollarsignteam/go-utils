package utils

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

func Min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func Max[T constraints.Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
}

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

func RandomInt64(min, max int64) int64 {
	if min > max {
		min, max = max, min
	}
	nBig, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
	return nBig.Int64() + min
}

func RandomFloat64(min, max float64) float64 {
	if min > max {
		min, max = max, min
	}
	nBig, _ := rand.Int(rand.Reader, big.NewInt(1<<62))
	return (float64(nBig.Int64())/float64(1<<62))*(max-min) + min
}

func ParseFloat64(s string) (float64, error) {
	s = strings.ReplaceAll(strings.TrimSpace(s), ",", "")
	return strconv.ParseFloat(s, 64)
}
