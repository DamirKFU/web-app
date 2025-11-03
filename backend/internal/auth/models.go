package auth

import (
	"app/internal/core"
	"app/internal/users"
)

type Session struct {
	core.AbstractModel
	UserID uint       `gorm:"not null;index"`
	User   users.User `gorm:"constraint:OnDelete:CASCADE;"`
}
