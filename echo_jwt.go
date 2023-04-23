package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// EchoJWT utility instance
var EchoJWT EchoJWTUtil

// EchoJWTUtil is a utility struct that provides methods
// for working with JWT tokens in the context of the Echo web framework
type EchoJWTUtil struct {
	config        EchoJWTConfig  // The configuration for EchoJWTUtil
	echoJWTConfig echojwt.Config // The configuration for the echojwt library
}

// JWTToken is a helper struct for returning signed JWT tokens
type JWTToken struct {
	SignedString string               // The signed token as a string
	Claims       jwt.RegisteredClaims // The claims included in the token
}

// EchoJWTConfig is the configuration struct for EchoJWTUtil
type EchoJWTConfig struct {
	SigningKey        string                                       // The signing key used to sign JWT tokens
	ExpiresTTL        time.Duration                                // The duration until which the token should be valid
	BeforeSuccessFunc func(token *jwt.Token, c echo.Context) error // A callback function to execute before a successful authentication
}

// New creates and returns a new instance of EchoJWTUtil
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

// CreateToken creates and returns a new JWTToken
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

// KeyFunc is a helper function used by ParseToken
// to extract the signing key from the EchoJWTConfig object
func (eJWT EchoJWTUtil) KeyFunc(token *jwt.Token) (interface{}, error) {
	if token.Method.Alg() != jwt.SigningMethodHS256.Name {
		return nil, fmt.Errorf("unexpected jwt signing method=%v", token.Header["alg"])
	}
	return []byte(eJWT.config.SigningKey), nil
}

// ParseToken is a helper function used to parse
// and validate JWT tokens using the echo-jwt library
func (eJWT EchoJWTUtil) ParseToken(signedToken string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(signedToken, &jwt.RegisteredClaims{}, eJWT.KeyFunc)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// ParseTokenFunc is a callback function used to parse
// and validate JWT tokens within the context of the echo-jwt middleware
func (eJWT EchoJWTUtil) ParseTokenFunc(c echo.Context, auth string) (interface{}, error) {
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

// JWTAuth returns a new instance of the echo-jwt middleware,
// configured with the current EchoJWTConfig object
func (eJWT EchoJWTUtil) JWTAuth() echo.MiddlewareFunc {
	return echojwt.WithConfig(eJWT.echoJWTConfig)
}

// GetClaims retrieves and validates JWT claims.
// It takes a JWT token and returns the converted claims
func (eJWT EchoJWTUtil) GetClaims(token *jwt.Token) (*jwt.RegisteredClaims, error) {
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}
