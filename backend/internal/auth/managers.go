package auth

import (
	"app/internal/core"

	"github.com/gin-gonic/gin"
)

type SessionManager struct {
	server *core.Server
}

func NewSessionManager(server *core.Server) *SessionManager {
	return &SessionManager{
		server: server,
	}
}

func (m *SessionManager) Create(c *gin.Context, userID uint) (*Session, error) {
	session := &Session{UserID: userID}
	if err := m.server.GetDB(c).Create(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func (m *SessionManager) GetByID(c *gin.Context, id uint) (*Session, error) {
	var session Session
	if err := m.server.GetDB(c).First(&session, id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (m *SessionManager) Delete(c *gin.Context, id uint) error {
	return m.server.GetDB(c).Delete(&Session{}, id).Error
}

func (m *SessionManager) DeleteByUser(c *gin.Context, userID uint) error {
	return m.server.GetDB(c).Where("user_id = ?", userID).Delete(&Session{}).Error
}
