package auth

import (
	"app/internal/core"
	"app/internal/users"
)

type Session struct {
	core.AbstractModel
	UserID uint       `gorm:"not null;index" binding:"SessionsUserID"`
	User   users.User `gorm:"constraint:OnDelete:CASCADE;"`
}
