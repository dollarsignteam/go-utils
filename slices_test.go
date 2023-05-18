package utils

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func Test_Map(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	double := func(x int) int { return x * 2 }
	r := Map(s, double)
	assert.Equal(t, len(r), len(s))
	assert.Equal(t, []int{2, 4, 6, 8, 10}, r)
}

func Test_Reduce(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	f := func(prev int, curr int) int { return prev + curr }
	r := Reduce(s, f, 0)
	assert.Equal(t, 15, r)
}

func Test_Reduce_string(t *testing.T) {
	s := []string{"a", "b", "c", "d", "e"}
	f := func(prev string, curr string) string { return prev + curr }
	r := Reduce(s, f, "")
	assert.Equal(t, "abcde", r)
}

func Test_Reduce_Float64(t *testing.T) {
	s := []float64{1.1, 2.2, 3.3, 4.4, 5.5}
	f := func(prev float64, curr float64) float64 { return prev + curr }
	r := Reduce(s, f, 0.0)
	assert.Equal(t, 16.5, r)
}

func Test_Filter(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	f := func(x int) bool { return x%2 == 0 }
	r := Filter(s, f)

	assert.Equal(t, len(r), 2)
	assert.ElementsMatch(t, []int{2, 4}, r)
}

func Test_SimpleSort(t *testing.T) {
	s := []int{5, 4, 3, 2, 56, 54, 56, 56, 1, 6, 12, 11}
	r := SimpleSort(s)
	assert.Equal(t, len(r), len(s))
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 11, 12, 54, 56, 56, 56}, r)
}

func Test_SortFunc(t *testing.T) {
	s := []int{5, 4, 3, 2, 1, 1, 1, 0, 1234, 12, 432, 34}
	f := func(x int, y int) bool { return y < x }
	r := SortFunc(s, f)

	assert.Equal(t, len(r), len(s))
	assert.Equal(t, []int{0, 1, 1, 1, 2, 3, 4, 5, 12, 34, 432, 1234}, r)
}

func Benchmark_SortFunc(b *testing.B) {
	s := getLargeSlice()
	f := func(x int, y int) bool { return x > y }
	for i := 0; i < b.N; i++ {
		SortFunc(s, f)
	}
}

func Benchmark_SimpleSort(b *testing.B) {
	s := getLargeSlice()
	for i := 0; i < b.N; i++ {
		SimpleSort(s)
	}
}

func getLargeSlice() []int {
	s := make([]int, 500)
	for i := 0; i < 500; i++ {
		s[i] = rand.Int()
	}
	return s
}
