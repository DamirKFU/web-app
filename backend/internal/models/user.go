package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string `gorm:"unique;not null;size:42"`
	Password  string `gorm:"not null"`
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
