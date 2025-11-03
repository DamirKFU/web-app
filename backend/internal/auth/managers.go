package auth

import (
	"app/internal/core"

	"gorm.io/gorm"
)

type SessionManager struct {
	server *core.Server
	db     *gorm.DB
}

func NewSessionManager(server *core.Server) *SessionManager {
	return &SessionManager{
		server: server,
		db:     server.DB,
	}
}

func (m *SessionManager) Create(userID uint) (*Session, error) {
	session := &Session{UserID: userID}
	if err := m.db.Create(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func (m *SessionManager) GetByID(id uint) (*Session, error) {
	var session Session
	if err := m.db.First(&session, id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (m *SessionManager) Delete(id uint) error {
	return m.db.Delete(&Session{}, id).Error
}

func (m *SessionManager) DeleteByUser(userID uint) error {
	return m.db.Where("user_id = ?", userID).Delete(&Session{}).Error
}
