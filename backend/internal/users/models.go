package users

import (
	"app/internal/core"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	core.AbstractModel
	Username    string `gorm:"type:varchar(42);unique;not null" json:"username"`
	Password    string `gorm:"type:varchar(128);unique;not null" json:"password"`
	IsSuperUser bool   `gorm:"default:false;not null" json:"is_super_user"`
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
