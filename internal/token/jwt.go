package token

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"github.com/web-programming-fall-2022/airline-auth/internal/storage"
	"strconv"
	"time"
)

type JWTManager struct {
	secret  string
	Storage *storage.Storage
	RDB     *redis.Client
}

func NewJWTManager(secret string, store *storage.Storage, rdb *redis.Client) *JWTManager {
	return &JWTManager{
		secret:  secret,
		Storage: store,
		RDB:     rdb,
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

	tokenString, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (m *JWTManager) Validate(ctx context.Context, tokenString string) (map[string]string, error) {
	err := m.CheckUnauthorizedToken(ctx, tokenString)
	if err != nil {
		return nil, err
	}
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

func (m *JWTManager) InvalidateToken(ctx context.Context, tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return m.secret, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token")
	}
	userId, _ := strconv.Atoi(claims["user_id"].(string))
	err = m.Storage.CreateUnauthorizedToken(&storage.UnauthorizedToken{
		UserID:     uint(userId),
		Token:      tokenString,
		Expiration: claims["exp"].(time.Time),
	})
	m.RDB.SetEx(ctx, tokenString, "true", time.Until(claims["exp"].(time.Time)))
	if err != nil {
		return err
	}
	return nil
}

func (m *JWTManager) CheckUnauthorizedToken(ctx context.Context, tokenString string) error {
	resp := m.RDB.Get(ctx, tokenString)
	if resp.Err() != nil {
		return nil
	}
	return errors.New("token is unauthorized")
}
