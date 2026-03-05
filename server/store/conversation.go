package store

import (
	"encoding/json"
	"errors"
	"openhands-go/server/models"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

type ConversationStore struct {
	mu            sync.RWMutex
	conversations map[string]models.ConversationInfo
	filePath      string
}

func NewConversationStore(filePath string) *ConversationStore {
	store := &ConversationStore{
		conversations: make(map[string]models.ConversationInfo),
		filePath:      filePath,
	}
	store.load()
	return store
}

func (s *ConversationStore) load() {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return
	}
	json.Unmarshal(data, &s.conversations)
}

func (s *ConversationStore) save() error {
	data, err := json.MarshalIndent(s.conversations, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0644)
}

func (s *ConversationStore) ListConversations() []models.ConversationInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conversations := make([]models.ConversationInfo, 0, len(s.conversations))
	for _, c := range s.conversations {
		conversations = append(conversations, c)
	}
	return conversations
}

func (s *ConversationStore) GetConversation(id string) (models.ConversationInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, ok := s.conversations[id]
	if !ok {
		return models.ConversationInfo{}, errors.New("conversation not found")
	}
	return c, nil
}

func (s *ConversationStore) CreateConversation(req models.InitSessionRequest) (models.ConversationInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	now := time.Now()

	conversation := models.ConversationInfo{
		ConversationID:     id,
		Title:              "New Conversation", // Default title logic to be implemented
		CreatedAt:          now,
		LastUpdatedAt:      now,
		Status:             models.ConversationStatusStopped, // Initially stopped
		SelectedRepository: req.Repository,
		SelectedBranch:     req.SelectedBranch,
		Trigger:            req.Trigger,
	}

	s.conversations[id] = conversation
	if err := s.save(); err != nil {
		return models.ConversationInfo{}, err
	}
	return conversation, nil
}

func (s *ConversationStore) DeleteConversation(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.conversations[id]; !ok {
		return errors.New("conversation not found")
	}

	delete(s.conversations, id)
	return s.save()
}

func (s *ConversationStore) UpdateConversation(id string, title string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, ok := s.conversations[id]
	if !ok {
		return errors.New("conversation not found")
	}

	c.Title = title
	c.LastUpdatedAt = time.Now()
	s.conversations[id] = c

	return s.save()
}
