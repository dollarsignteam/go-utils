package utils

import "testing"

func Test_Map(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	double := func(x int) int { return x * 2 }
	r := Map(s, double)
	if len(r) != len(s) {
		t.Errorf("Map: expected %v, got %v", len(s), len(r))
	}
	for i, v := range r {
		if v != s[i]*2 {
			t.Errorf("Map: expected %v, got %v", s[i]*2, v)
		}
	}
}

func Test_Reduce(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	f := func(prev int, curr int) int { return prev + curr }
	r := Reduce(s, f, 0)
	if r != 15 {
		t.Errorf("Reduce: expected %v, got %v", 15, r)
	}
}

func Test_Filter(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	f := func(x int) bool { return x%2 == 0 }
	r := Filter(s, f)
	if len(r) != 2 {
		t.Errorf("Filter: expected %v, got %v", 2, len(r))
	}
	for _, v := range r {
		if v%2 != 0 {
			t.Errorf("Filter: expected %v, got %v", true, false)
		}
	}
}
