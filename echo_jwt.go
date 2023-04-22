package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

var EchoJWT EchoJWTUtil

type EchoJWTUtil struct {
	config        EchoJWTConfig
	echoJWTConfig echojwt.Config
}

type JWTToken struct {
	SignedString string
	Claims       jwt.RegisteredClaims
}

type EchoJWTConfig struct {
	SigningKey        string
	ExpiresTTL        time.Duration
	BeforeSuccessFunc func(token *jwt.Token, c echo.Context) error
}

func (EchoJWTUtil) New(config EchoJWTConfig) *EchoJWTUtil {
	echoJWTUtil := &EchoJWTUtil{
		config: config,
		echoJWTConfig: echojwt.Config{
			SigningKey: []byte(config.SigningKey),
		},
	}
	echoJWTUtil.echoJWTConfig.ParseTokenFunc = echoJWTUtil.ParseTokenFunc
	return echoJWTUtil
}

func (eJWT EchoJWTUtil) CreateToken(claims jwt.RegisteredClaims) JWTToken {
	if claims.ID == "" {
		claims.ID = String.UUID()
	}
	if claims.IssuedAt == nil {
		claims.IssuedAt = jwt.NewNumericDate(time.Now())
	}
	if claims.ExpiresAt == nil {
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(eJWT.config.ExpiresTTL))
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(eJWT.config.SigningKey))
	return JWTToken{
		SignedString: signedToken,
		Claims:       claims,
	}
}

func (eJWT EchoJWTUtil) KeyFunc(token *jwt.Token) (any, error) {
	if token.Method.Alg() != jwt.SigningMethodHS256.Name {
		return nil, fmt.Errorf("unexpected jwt signing method=%v", token.Header["alg"])
	}
	return []byte(eJWT.config.SigningKey), nil
}

func (eJWT EchoJWTUtil) ParseToken(signedToken string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(signedToken, &jwt.RegisteredClaims{}, eJWT.KeyFunc)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (eJWT EchoJWTUtil) ParseTokenFunc(c echo.Context, auth string) (any, error) {
	token, err := eJWT.ParseToken(auth)
	if err != nil {
		return nil, err
	}
	if eJWT.config.BeforeSuccessFunc != nil {
		if err := eJWT.config.BeforeSuccessFunc(token, c); err != nil {
			return nil, err
		}
	}
	return token, nil
}

func (eJWT EchoJWTUtil) JWTAuth() echo.MiddlewareFunc {
	return echojwt.WithConfig(eJWT.echoJWTConfig)
}
