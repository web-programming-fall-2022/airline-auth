package storage

import (
	"gorm.io/gorm"
	"time"
)

type UnauthorizedToken struct {
	gorm.Model
	User       AuthUser
	Token      string
	Expiration time.Time
}
