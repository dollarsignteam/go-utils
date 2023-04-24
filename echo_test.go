package utils_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/dollarsignteam/go-utils"
)

func TestDefaultRootHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := utils.Echo.DefaultRootHandler(c)
	expected := `{"message":"200 OK"}`
	actual := strings.TrimRight(rec.Body.String(), "\n")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expected, actual)
	assert.NoError(t, err)
}

func TestNoContentHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := utils.Echo.NoContentHandler(c)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.NoError(t, err)
}

func TestEchoBinderWithValidation_Bind(t *testing.T) {
	type TestRequest struct {
		Field1 string `json:"field1" validate:"required"`
		Field2 int    `json:"field2" validate:"gte=0"`
	}
	tests := []struct {
		body     string
		expected error
	}{
		{
			body:     `{"field1":"value1", "field2":1}`,
			expected: nil,
		},
		{
			body:     `{"field2":1}`,
			expected: errors.New("Validation failed for 'field1'"),
		},
		{
			body:     `[]`,
			expected: errors.New("Unmarshal type error: expected=utils_test.TestRequest, got=array, field=, offset=1"),
		},
	}
	e := echo.New()
	binder := utils.EchoBinder
	for _, test := range tests {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, httptest.NewRecorder())
		testRequest := new(TestRequest)
		err := binder.Bind(testRequest, c)
		if err != nil {
			assert.EqualError(t, err, test.expected.Error())
		} else {
			assert.Equal(t, test.expected, err)
		}
	}
}

func TestEchoBinderWithValidation_BindBody(t *testing.T) {
	type TestRequest struct {
		Field1 string `json:"field1" validate:"required"`
		Field2 int    `json:"field2" validate:"gte=0"`
	}
	tests := []struct {
		body     string
		expected error
	}{
		{
			body:     `{"field1":"value1", "field2":1}`,
			expected: nil,
		},
		{
			body:     `[]`,
			expected: errors.New("Unmarshal type error: expected=utils_test.TestRequest, got=array, field=, offset=1"),
		},
	}
	e := echo.New()
	binder := utils.EchoBinder
	for _, test := range tests {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, httptest.NewRecorder())
		testRequest := new(TestRequest)
		err := binder.BindBody(c, testRequest)
		if err != nil {
			assert.EqualError(t, err, test.expected.Error())
		} else {
			assert.Equal(t, test.expected, err)
		}
	}
}

func TestEchoBinderWithValidation_BindHeaders(t *testing.T) {
	type TestRequest struct {
		Field1 string `header:"field1" validate:"required"`
		Field2 int    `header:"field2" validate:"gte=0"`
	}
	tests := []struct {
		name           string
		headers        map[string]string
		expectedResult interface{}
		expectedError  error
	}{
		{
			name: "Successful binding and validation",
			headers: map[string]string{
				"Field1": "hello",
				"Field2": "3",
			},
			expectedResult: &TestRequest{
				Field1: "hello",
				Field2: 3,
			},
			expectedError: nil,
		},
		{
			name: "Missing required header field",
			headers: map[string]string{
				"Field2": "5",
			},
			expectedResult: nil,
			expectedError:  errors.New("Validation failed for 'Field1'"),
		},
	}
	e := echo.New()
	binder := utils.EchoBinder
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(echo.GET, "/", nil)
			for k, v := range test.headers {
				req.Header.Add(k, v)
			}
			c := e.NewContext(req, httptest.NewRecorder())
			testRequest := new(TestRequest)
			err := binder.BindHeaders(c, testRequest)
			if test.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, test.expectedResult, testRequest)
			}
		})
	}
	req := httptest.NewRequest(echo.GET, "/", nil)
	c := e.NewContext(req, httptest.NewRecorder())
	err := binder.BindHeaders(c, nil)
	assert.EqualError(t, err, "")
}

func TestEchoValidator_Validate(t *testing.T) {
	validator := &utils.EchoValidator{}
	var data struct {
		Name string `json:"name" validate:"required"`
	}
	err := validator.Validate(&data)
	assert.EqualError(t, err, "Validation failed for 'name'")
	data.Name = "John"
	err = validator.Validate(&data)
	assert.NoError(t, err)
}

func TestNew(t *testing.T) {
	e := utils.Echo.New()
	assert.IsType(t, &echo.Echo{}, e)
	assert.True(t, e.HidePort)
	assert.True(t, e.HideBanner)
	assert.IsType(t, &utils.EchoValidator{}, e.Validator)
	assert.IsType(t, &utils.EchoBinderWithValidation{}, e.Binder)
}
