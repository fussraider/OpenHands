package store

import (
	"errors"
	"log/slog"
	"openhands-go/server/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConversationStore struct {
	db *gorm.DB
}

func NewConversationStore() *ConversationStore {
	return &ConversationStore{
		db: DB,
	}
}

func (s *ConversationStore) ListConversations() []models.ConversationInfo {
	var conversations []models.ConversationInfo
	s.db.Find(&conversations)
	return conversations
}

func (s *ConversationStore) GetConversation(id string) (models.ConversationInfo, error) {
	var c models.ConversationInfo
	result := s.db.First(&c, "conversation_id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c, errors.New("conversation not found")
		}
		return c, result.Error
	}
	return c, nil
}

func (s *ConversationStore) CreateConversation(req models.InitSessionRequest) (models.ConversationInfo, error) {
	var count int64
	s.db.Model(&models.ConversationInfo{}).Count(&count)
	if count >= 50 {
		slog.Debug("closing_from_too_many_sessions", "warning", "max limit reached")
	}

	id := uuid.New().String()
	now := time.Now()

	conversation := models.ConversationInfo{
		ConversationID:     id,
		Title:              "New Conversation",
		CreatedAt:          now,
		LastUpdatedAt:      now,
		Status:             models.ConversationStatusStopped,
		SelectedRepository: req.Repository,
		SelectedBranch:     req.SelectedBranch,
		Trigger:            req.Trigger,
	}

	result := s.db.Create(&conversation)
	if result.Error != nil {
		return models.ConversationInfo{}, result.Error
	}
	return conversation, nil
}

func (s *ConversationStore) DeleteConversation(id string) error {
	result := s.db.Delete(&models.ConversationInfo{}, "conversation_id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("conversation not found")
	}
	return nil
}

func (s *ConversationStore) UpdateConversation(id string, title string) error {
	result := s.db.Model(&models.ConversationInfo{}).Where("conversation_id = ?", id).Updates(models.ConversationInfo{
		Title:         title,
		LastUpdatedAt: time.Now(),
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("conversation not found")
	}
	return nil
}

func (s *ConversationStore) SetConversationStatus(id string, status models.ConversationStatus) error {
	result := s.db.Model(&models.ConversationInfo{}).Where("conversation_id = ?", id).Updates(map[string]interface{}{
		"status":          status,
		"last_updated_at": time.Now(),
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("conversation not found")
	}
	return nil
}
