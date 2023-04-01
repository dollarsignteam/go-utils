package utils_test

import (
	"fmt"

	"github.com/dollarsignteam/go-utils"
)

func ExamplePointerOf() {
	x := 20
	p := utils.PointerOf(x)
	fmt.Println(*p)
	// Output: 20
}

func ExamplePackageName() {
	fmt.Println(utils.PackageName())
	// Output: go-utils_test
}

func ExampleUniqueOf() {
	s := []int{1, 2, 2, 3, 4, 4, 5, 5, 5}
	u := utils.UniqueOf(s)
	fmt.Println(u)
	// Output: [1 2 3 4 5]
}

func ExampleValueOf() {
	x := 20
	p := utils.PointerOf(x)
	v := utils.ValueOf(p)
	fmt.Println(v)
	// Output: 20
}

func ExampleIsArrayOrSlice() {
	s := []int{1, 2, 3, 4, 5}
	fmt.Println(utils.IsArrayOrSlice(s))
	// Output: true
}

func ExampleBoolToInt() {
	b := true
	i := utils.BoolToInt(b)
	fmt.Println(i)
	// Output: 1
}

func ExampleIntToBool() {
	i := 0
	b := utils.IntToBool(i)
	fmt.Println(b)
	// Output: false
}
