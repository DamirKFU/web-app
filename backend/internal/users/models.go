package users

import (
	"app/internal/core"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	core.AbstractModel
	Username       string `gorm:"type:varchar(42);unique;not null" json:"username"`
	Email          string `gorm:"type:varchar(128);unique;not null" json:"email"`
	NormalizeEmail string `gorm:"type:varchar(128);unique;not null" json:"normalize_email"`
	Password       string `gorm:"type:varchar(128);not null" json:"password"`
	IsSuperUser    bool   `gorm:"default:false;not null" json:"is_super_user"`
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

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.NormalizeEmail = NormalizeEmail(u.Email)
	return nil
}
