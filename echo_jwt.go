package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

var EchoJWT EchoJWTUtil

type EchoJWTUtil struct {
	signingKey    []byte
	expiresTTL    time.Duration
	echoJWTConfig echojwt.Config
}

type EchoJWTConfig struct {
	SigningKey string
	ExpiresTTL time.Duration
}

type JWTToken struct {
	ID     string
	Token  string
	Claims jwt.RegisteredClaims
}

func (EchoJWTUtil) New(config EchoJWTConfig) *EchoJWTUtil {
	signingKey := []byte(config.SigningKey)
	return &EchoJWTUtil{
		signingKey: signingKey,
		expiresTTL: config.ExpiresTTL,
		echoJWTConfig: echojwt.Config{
			SigningKey: signingKey,
		},
	}
}

func (eJWT EchoJWTUtil) JWTAuth() echo.MiddlewareFunc {
	return echojwt.WithConfig(eJWT.echoJWTConfig)
}

func (eJWT EchoJWTUtil) CreateToken(claims jwt.RegisteredClaims) (JWTToken, error) {
	if claims.ID == "" {
		claims.ID = String.UUID()
	}
	timeNow := time.Now()
	if claims.IssuedAt == nil {
		claims.IssuedAt = jwt.NewNumericDate(timeNow)
	}
	if claims.ExpiresAt == nil {
		claims.ExpiresAt = jwt.NewNumericDate(timeNow.Add(eJWT.expiresTTL))
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString(eJWT.signingKey)
	if err != nil {
		return JWTToken{}, err
	}
	return JWTToken{
		ID:     claims.ID,
		Token:  token,
		Claims: claims,
	}, nil
}

func (eJWT EchoJWTUtil) ParseToken(tokenString string) (jwt.MapClaims, error) {
	// Parse the token and check for validity and expiration
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return eJWT.signingKey, nil
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
