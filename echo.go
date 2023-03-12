package utils

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Echo utility instance
var Echo EchoUtil

// EchoUtil is a utility struct for working with Echo instances
type EchoUtil struct{}

// EchoValidator is a struct that implements the echo.Validator interface.
type EchoValidator struct{}

// Validate is a method that validates the given struct using the ValidateStruct
// function and returns an error if validation fails.
func (EchoValidator) Validate(i any) error {
	return ValidateStruct(i)
}

// EchoBinderWithValidation is a struct that implements the echo.Binder interface
// with added validation functionality.
type EchoBinderWithValidation struct {
	echo.DefaultBinder
}

// Bind is a method that binds the request data to the given struct,
// validates it using the ValidateStruct function,
// and returns an error if binding or validation fails.
func (b *EchoBinderWithValidation) Bind(i any, c echo.Context) error {
	if err := b.DefaultBinder.Bind(i, c); err != nil {
		return errors.New(err.(*echo.HTTPError).Message.(string))
	}
	return ValidateStruct(i)
}

// DefaultRootHandler handles requests to the root endpoint
func (EchoUtil) DefaultRootHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "200 OK"})
}

// New creates a new instance of the Echo framework
func (EchoUtil) New() *echo.Echo {
	e := echo.New()
	e.HidePort = true
	e.HideBanner = true
	e.Validator = &EchoValidator{}
	e.Binder = &EchoBinderWithValidation{}
	e.Pre(middleware.RemoveTrailingSlash())
	e.GET("/", Echo.DefaultRootHandler)
	return e
}
