package utils_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/dollarsignteam/go-utils"
)

var testEchoJWTConfig = utils.EchoJWTConfig{
	SigningKey: "my-secret-key",
	ExpiresTTL: time.Hour,
}
var testClaims = jwt.RegisteredClaims{
	Subject: "test-user",
}

func TestCreateToken(t *testing.T) {
	echoJWT := utils.EchoJWT.New(testEchoJWTConfig)
	token := echoJWT.CreateToken(testClaims)
	assert.NotEmpty(t, token.SignedString)
	assert.Equal(t, testClaims.Subject, token.Claims.Subject)
}

func TestKeyFunc(t *testing.T) {
	echoJWT := utils.EchoJWT.New(testEchoJWTConfig)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)
	result, err := echoJWT.KeyFunc(token)
	assert.Nil(t, err)
	assert.Equal(t, []byte(testEchoJWTConfig.SigningKey), result)
}

func TestKeyFunc_Error(t *testing.T) {
	echoJWT := utils.EchoJWT.New(testEchoJWTConfig)
	token := jwt.NewWithClaims(jwt.SigningMethodNone, testClaims)
	result, err := echoJWT.KeyFunc(token)
	assert.Nil(t, result)
	assert.EqualError(t, err, "unexpected jwt signing method=none")
}

func TestParseToken(t *testing.T) {
	echoJWT := utils.EchoJWT.New(testEchoJWTConfig)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)
	signedToken, _ := token.SignedString([]byte(testEchoJWTConfig.SigningKey))
	result, err := echoJWT.ParseToken(signedToken)
	assert.True(t, result.Valid)
	assert.NoError(t, err)
}

func TestParseToken_InvalidClaims(t *testing.T) {
	echoJWT := utils.EchoJWT.New(testEchoJWTConfig)
	claims := jwt.MapClaims{"jti": true}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(testEchoJWTConfig.SigningKey))
	result, err := echoJWT.ParseToken(signedToken)
	assert.Nil(t, result)
	assert.EqualError(t, err, "json: cannot unmarshal bool into Go struct field RegisteredClaims.jti of type string")
}

func TestParseToken_InvalidSigningKey(t *testing.T) {
	echoJWT := utils.EchoJWT.New(testEchoJWTConfig)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)
	signedToken, _ := token.SignedString([]byte("wrong-secret-key"))
	_, err := echoJWT.ParseToken(signedToken)
	assert.EqualError(t, err, "signature is invalid")
}

func TestParseTokenFunc(t *testing.T) {
	echoJWT := utils.EchoJWT.New(testEchoJWTConfig)
	_, err := echoJWT.ParseTokenFunc(nil, "")
	assert.EqualError(t, err, "token contains an invalid number of segments")
}

func TestParseTokenFunc_InvalidToken(t *testing.T) {
	echoJWT := utils.EchoJWT.New(testEchoJWTConfig)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)
	signedToken, _ := token.SignedString([]byte(testEchoJWTConfig.SigningKey))
	result, err := echoJWT.ParseTokenFunc(nil, signedToken)
	assert.Nil(t, err)
	assert.IsType(t, &jwt.Token{}, result)
}

func TestParseTokenFunc_BeforeSuccessFunc_Error(t *testing.T) {
	var beforeSuccessCalled bool
	config := utils.EchoJWTConfig{
		SigningKey: testEchoJWTConfig.SigningKey,
		BeforeSuccessFunc: func(token *jwt.Token, c echo.Context) error {
			beforeSuccessCalled = true
			return errors.New("mock error")
		},
	}
	echoJWT := utils.EchoJWT.New(config)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)
	signedToken, _ := token.SignedString([]byte(testEchoJWTConfig.SigningKey))
	result, err := echoJWT.ParseTokenFunc(nil, signedToken)
	assert.Nil(t, result)
	assert.EqualError(t, err, "mock error")
	assert.True(t, beforeSuccessCalled)
}

func TestJWTAuth(t *testing.T) {
	echoJWT := utils.EchoJWT.New(testEchoJWTConfig)
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		token := c.Get("user").(*jwt.Token)
		claims := token.Claims.(*jwt.RegisteredClaims)
		return c.String(http.StatusOK, claims.Subject)
	})
	e.Use(echoJWT.JWTAuth())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	token := echoJWT.CreateToken(testClaims)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token.SignedString))
	e.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, testClaims.Subject, res.Body.String())
}
