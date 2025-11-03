package users

import (
	"app/internal/core"
)

type UserManager struct {
	server *core.Server
}

func NewUserManager(server *core.Server) *UserManager {
	return &UserManager{server: server}
}

func (m *UserManager) GetByUsername(username string) (*User, error) {
	var user User
	if err := m.server.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UserManager) GetByEmail(email string) (*User, error) {
	var user User
	if err := m.server.DB.Where("normalize_email = ?", NormalizeEmail(email)).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UserManager) GetByID(id uint) (*User, error) {
	var user User
	if err := m.server.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UserManager) Create(user *User) error {
	return m.server.DB.Create(user).Error
}

func (m *UserManager) Save(user *User) error {
	return m.server.DB.Save(user).Error
}
