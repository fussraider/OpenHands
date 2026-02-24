package session

import (
	"errors"
	"sync"
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
	mu       sync.RWMutex
	sessions = make(map[string]Session)
)

func CreateSession(userID string) (string, error) {
	mu.Lock()
	defer mu.Unlock()
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
	mu.RLock()
	session, ok := sessions[sessionID]
	mu.RUnlock()

	if !ok {
		return Session{}, errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		mu.Lock()
		delete(sessions, sessionID)
		mu.Unlock()
		return Session{}, errors.New("session expired")
	}
	return session, nil
}
