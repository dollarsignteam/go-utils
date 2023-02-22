# Utilities for the Go

Utilities for the Golang package. The utils package is a generic package that provides utility functions and structures that can be used across different parts of a program. The package contains several functions and structures that can be used to simplify common tasks.

## Installation

```shell
go get github.com/dollarsignteam/go-utils
```

## Usage

Here is an example of how to use the RandomInt64 function to generate a random integer between 1 and 100:

```go
package main

import (
    "fmt"
    "github.com/dollarsignteam/go-utils"
)

func main() {
    randInt := utils.RandomInt64(1, 100)
    fmt.Printf("Random integer: %d\n", randInt)
}
```

## Author

Dollarsign

## License

Licensed under the MIT License - see the [LICENSE][1] file for details.

[1]: https://github.com/dollarsignteam/go-logger/blob/main/LICENSE
