package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var EchoJWT EchoJWTUtil

type EchoJWTUtil struct {
	config EchoJWTConfig
}

type EchoJWTConfig struct {
	SigningKey     string
	ExpiresTTL     time.Duration
	ParseTokenFunc func(c echo.Context, auth string) (any, error)
}

func (EchoJWTUtil) New(config EchoJWTConfig) *EchoJWTUtil {
	return &EchoJWTUtil{
		config: config,
	}
}

func (eJWT EchoJWTUtil) JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the token from the Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.String(http.StatusUnauthorized, "missing Authorization header")
			}

			// Split the token from the "Bearer " prefix
			auth := strings.Split(authHeader, "Bearer ")
			if len(auth) != 2 {
				return c.String(http.StatusUnauthorized, "invalid Authorization header format")
			}

			// Parse the token and check for validity and expiration
			token, err := eJWT.ParseToken(auth[1])
			if err != nil {
				return c.String(http.StatusUnauthorized, fmt.Sprintf("invalid token: %v", err))
			}

			// Set the user information in the echo context and continue the request
			c.Set("user", token)
			return next(c)
		}
	}
}

func (eJWT EchoJWTUtil) CreateToken(claims jwt.MapClaims) (string, error) {
	// Set the token expiration time
	expirationTime := time.Now().Add(eJWT.config.ExpiresTTL).Unix()

	// Add the expiration to the claims
	claims["exp"] = expirationTime

	// Create the JWT token with the claims and signing key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(eJWT.config.SigningKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (eJWT EchoJWTUtil) ParseToken(tokenString string) (jwt.MapClaims, error) {
	// Parse the token and check for validity and expiration
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(eJWT.config.SigningKey), nil
	})
	if err != nil {
		return nil, err
	}

	// Check if the token is still valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Get the claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

func (eJWT EchoJWTUtil) GetSession(c echo.Context) (jwt.MapClaims, error) {
	// Get the token from the Authorization header
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("missing Authorization header")
	}

	// Split the token from the "Bearer " prefix
	auth := strings.Split(authHeader, "Bearer ")
	if len(auth) != 2 {
		return nil, errors.New("invalid Authorization header format")
	}

	// Parse the token and return the claims
	claims, err := eJWT.ParseToken(auth[1])
	if err != nil {
		return nil, fmt.Errorf("unable to parse token: %v", err)
	}

	return claims, nil
}
