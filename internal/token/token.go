package token

import (
	"time"
)

type Manager interface {
	Generate(claims map[string]string, expiration time.Time) (string, error)
	Validate(token string) (map[string]string, error)
}
