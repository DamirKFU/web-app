package users

import (
	"app/internal/core"

	"github.com/gin-gonic/gin"
)

type UserManager struct {
	server *core.Server
}

func NewUserManager(server *core.Server) *UserManager {
	return &UserManager{server: server}
}

func (m *UserManager) GetByUsername(c *gin.Context, username string) (*User, error) {
	var user User
	if err := m.server.GetDB(c).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UserManager) GetByEmail(c *gin.Context, email string) (*User, error) {
	var user User
	if err := m.server.GetDB(c).Where("normalize_email = ?", NormalizeEmail(email)).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UserManager) GetByID(c *gin.Context, id uint) (*User, error) {
	var user User
	if err := m.server.GetDB(c).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UserManager) Create(c *gin.Context, user *User) error {
	return m.server.GetDB(c).Create(user).Error
}

func (m *UserManager) Save(c *gin.Context, user *User) error {
	return m.server.GetDB(c).Save(user).Error
}
