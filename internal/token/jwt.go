package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JWTManager struct {
	secret string
}

func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{
		secret: secret,
	}
}

func (m *JWTManager) Generate(claims map[string]string, expiration time.Time) (string, error) {
	mapClaims := jwt.MapClaims{
		"exp": expiration.Unix(),
	}
	for key, value := range claims {
		mapClaims[key] = value
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)

	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (m *JWTManager) Validate(tokenString string) (map[string]string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}
	result := make(map[string]string)
	for key, value := range claims {
		result[key] = value.(string)
	}
	return result, nil
}
