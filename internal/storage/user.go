package storage

import (
	"errors"
	"gorm.io/gorm"
)

type AuthUser struct {
	gorm.Model
	Email        string `gorm:"index_unique"`
	PhoneNumber  string
	Gender       string
	FirstName    string
	LastName     string
	PasswordHash string
}

func (storage *Storage) CreateUser(user *AuthUser) error {
	if err := storage.DB.Create(user).Error; err != nil {
		return errors.New("couldn't create user in postgres storage")
	}
	return nil
}

func (storage *Storage) GetUserByID(id uint) (*AuthUser, error) {
	user := AuthUser{}
	storage.DB.First(&user, id)
	if user.ID == 0 {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (storage *Storage) GetUserByEmail(email string) (*AuthUser, error) {
	user := AuthUser{Email: email}
	storage.DB.Where(&user).First(&user)
	if user.ID == 0 {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
