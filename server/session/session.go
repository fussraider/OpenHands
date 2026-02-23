package session

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        string
	UserID    string
	CreatedAt time.Time
	ExpiresAt time.Time
}

var (
	sessions = make(map[string]Session)
)

func CreateSession(userID string) (string, error) {
	sessionID := uuid.New().String()
	sessions[sessionID] = Session{
		ID:        sessionID,
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	return sessionID, nil
}

func GetSession(sessionID string) (Session, error) {
	session, ok := sessions[sessionID]
	if !ok {
		return Session{}, errors.New("session not found")
	}
	if time.Now().After(session.ExpiresAt) {
		delete(sessions, sessionID)
		return Session{}, errors.New("session expired")
	}
	return session, nil
}
