package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

func (u *User) SetPassword(password string, secretKey string) error {
	peppered := password + secretKey
	bytes, err := bcrypt.GenerateFromPassword([]byte(peppered), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

func (u *User) CheckPassword(password string, secretKey string) bool {
	peppered := password + secretKey
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(peppered)) == nil
}
