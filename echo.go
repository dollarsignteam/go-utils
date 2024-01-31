package utils

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Echo utility instance
var Echo EchoUtil

// EchoBinder utility instance
var EchoBinder EchoBinderWithValidation

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

// validateWithErrorHandling validates a given struct and error handler
func (b *EchoBinderWithValidation) validateWithErrorHandling(i any, err error) error {
	if err != nil {
		return errors.New(err.(*echo.HTTPError).Message.(string))
	}
	return ValidateStruct(i)
}

// Bind binds request data, validates it using ValidateStruct(),
// and returns an error if binding or validation fails.
func (b *EchoBinderWithValidation) Bind(i any, c echo.Context) error {
	err := b.DefaultBinder.Bind(i, c)
	return b.validateWithErrorHandling(i, err)
}

// BindBody binds body data, validates it using ValidateStruct(),
// and returns an error if binding or validation fails.
func (b *EchoBinderWithValidation) BindBody(c echo.Context, i any) error {
	err := b.DefaultBinder.BindBody(c, i)
	return b.validateWithErrorHandling(i, err)
}

// BindHeaders binds headers data, validates it using ValidateStruct(),
// and returns an error if binding or validation fails.
func (b *EchoBinderWithValidation) BindHeaders(c echo.Context, i any) error {
	err := b.DefaultBinder.BindHeaders(c, i)
	return b.validateWithErrorHandling(i, err)
}

// BindPathParams binds path params, validates them using ValidateStruct(),
// and returns an error if binding or validation fails.
func (b *EchoBinderWithValidation) BindPathParams(c echo.Context, i any) error {
	err := b.DefaultBinder.BindPathParams(c, i)
	return b.validateWithErrorHandling(i, err)
}

// BindQueryParams binds query params, validates them using ValidateStruct(),
// and returns an error if binding or validation fails.
func (b *EchoBinderWithValidation) BindQueryParams(c echo.Context, i any) error {
	err := b.DefaultBinder.BindQueryParams(c, i)
	return b.validateWithErrorHandling(i, err)
}

// BindAll binds all request data, validates it using ValidateStruct(),
func (b *EchoBinderWithValidation) BindAll(i any, c echo.Context) error {
	if err := b.DefaultBinder.BindPathParams(c, i); err != nil {
		return b.validateWithErrorHandling(i, err)
	}
	if err := b.DefaultBinder.BindQueryParams(c, i); err != nil {
		return b.validateWithErrorHandling(i, err)
	}
	err := b.DefaultBinder.BindBody(c, i)
	return b.validateWithErrorHandling(i, err)
}

// DefaultRootHandler handles requests to the root endpoint
func (EchoUtil) DefaultRootHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "200 OK"})
}

// NoContentHandler handles return no content endpoint
func (EchoUtil) NoContentHandler(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// New creates a new instance of the Echo framework
func (EchoUtil) New() *echo.Echo {
	e := echo.New()
	e.HidePort = true
	e.HideBanner = true
	e.Validator = new(EchoValidator)
	e.Binder = &EchoBinder
	e.Pre(middleware.RemoveTrailingSlash())
	e.GET("/", Echo.DefaultRootHandler)
	e.GET("/favicon.ico", Echo.NoContentHandler)
	return e
}
