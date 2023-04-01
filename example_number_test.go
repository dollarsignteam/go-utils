package utils_test

import (
	"fmt"

	"github.com/dollarsignteam/go-utils"
)

func ExampleMin() {
	fmt.Println(utils.Min(4, 9))
	fmt.Println(utils.Min(9, 4))
	// Output: 4
	// 4
}

func ExampleMax() {
	fmt.Println(utils.Max(4, 9))
	fmt.Println(utils.Max(9, 4))
	// Output: 9
	// 9
}

func ExampleMinOf() {
	fmt.Println(utils.MinOf([]int{5, 2, 8, 1}))
	fmt.Println(utils.MinOf([]int{9, 10, 20}))
	fmt.Println(utils.MinOf([]int{}))
	// Output: 1
	// 9
	// 0
}

func ExampleMaxOf() {
	fmt.Println(utils.MaxOf([]int{5, 2, 8, 1}))
	fmt.Println(utils.MaxOf([]int{9, 10, 20}))
	fmt.Println(utils.MaxOf([]int{}))
	// Output: 8
	// 20
	// 0
}

func ExampleParseFloat64() {
	f, err := utils.ParseFloat64("3.14")
	if err == nil {
		fmt.Println(f)
	}
	// Output: 3.14
}

func ExampleRandomInt64() {
	fmt.Println(utils.RandomInt64(0, 10))
}

func ExampleRandomFloat64() {
	fmt.Println(utils.RandomFloat64(0, 10))
}
