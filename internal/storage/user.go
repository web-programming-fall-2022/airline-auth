package storage

import "gorm.io/gorm"

type AuthUser struct {
	gorm.Model
	Email        string
	PhoneNumber  string
	Gender       string
	FirstName    string
	LastName     string
	PasswordHash string
}
