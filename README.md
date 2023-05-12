# Utilities for the Go

Utilities for the Golang package. The utils package is a generic package that provides utility functions and structures that can be used across different parts of a program. The package contains several functions and structures that can be used to simplify common tasks.

## Installation

```shell
go get github.com/dollarsignteam/go-utils
```

## Usage

`PointerOf` returns a pointer to the input value. For example:

```go
x := 42
ptr := PointerOf(x)
fmt.Println(*ptr) // Output: 42
```

`PackageName` returns the name of the package that calls it. For example:

```go
fmt.Println(PackageName()) // Output: utils
```

`UniqueOf` removes duplicates from a slice of any type and returns a new slice containing only the unique elements. For example:

```go
input := []int{1, 2, 3, 2, 1}
unique := UniqueOf(input)
fmt.Println(unique) // Output: [1 2 3]
```

`ValueOf` takes a pointer to a value of any type and returns the value. For example:

```go
x := 42
ptr := &x
val := ValueOf(ptr)
fmt.Println(val) // Output: 42
```

`IsArrayOrSlice` takes a value of any type and returns a boolean indicating if it is a slice or an array. For example:

```go
arr := [3]int{1, 2, 3}
slice := []int{1, 2, 3}
fmt.Println(IsArrayOrSlice(arr)) // Output: true
fmt.Println(IsArrayOrSlice(slice)) // Output: true
fmt.Println(IsArrayOrSlice(x)) // Output: false
```

`BoolToInt` converts a boolean value to an integer (1 for true, 0 for false). For example:

```go
fmt.Println(BoolToInt(true)) // Output: 1
fmt.Println(BoolToInt(false)) // Output: 0
```

`IntToBool` converts an integer value to a boolean (true for non-zero values, false for zero). For example:

```go
fmt.Println(IntToBool(1)) // Output: true
fmt.Println(IntToBool(0)) // Output: false
```

For more information, check out the 📚 [documentation][2].

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## Author

Dollarsign

## License

Licensed under the MIT License - see the [LICENSE][1] file for details.

[1]: https://github.com/dollarsignteam/go-utils/blob/main/LICENSE
[2]: https://pkg.go.dev/github.com/dollarsignteam/go-utils
