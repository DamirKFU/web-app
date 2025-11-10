package users

import (
	"app/internal/core"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	core.AbstractModel
	Username       string `gorm:"type:varchar(42);unique;not null" binding:"UsersUsername"`
	Email          string `gorm:"type:varchar(128);unique;not null" binding:"UsersEmail"`
	NormalizeEmail string `gorm:"type:varchar(128);unique;not null" binding:"UsersEmail"`
	Password       string `gorm:"type:varchar(128);not null" binding:"UsersPassword"`
	IsSuperUser    bool   `gorm:"default:false;not null"`
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	u.NormalizeEmail = NormalizeEmail(u.Email)
	return core.ValidateStruct(u)
}

func (u *User) IsValid() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		return v.Struct(u)
	}
	return nil
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
