package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrCannotParsePayload = errors.New("cannot parse")
)

func CreateJWT(payload jwt.MapClaims, method jwt.SigningMethod, secret string) (string, error) {
	return jwt.NewWithClaims(method, payload).SignedString([]byte(secret))
}

func ValidateToken(tokenString, secret string) (map[string]any, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) { // todo: see doc of parse
		return []byte(secret), nil
	})
	if err != nil { // todo: jwt.WithValidMethods()
		return nil, err
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrCannotParsePayload
	}

	return payload, nil
}
