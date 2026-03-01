package models

import (
	"time"
)

type ConversationStatus string

const (
	ConversationStatusRunning  ConversationStatus = "running"
	ConversationStatusStopped  ConversationStatus = "stopped"
	ConversationStatusStarting ConversationStatus = "starting"
	ConversationStatusError    ConversationStatus = "error"
)

type ConversationInfo struct {
	ConversationID     string             `json:"conversation_id"`
	Title              string             `json:"title"`
	LastUpdatedAt      time.Time          `json:"last_updated_at"`
	CreatedAt          time.Time          `json:"created_at"`
	Status             ConversationStatus `json:"status"`
	SelectedRepository string             `json:"selected_repository,omitempty"`
	SelectedBranch     string             `json:"selected_branch,omitempty"`
	Trigger            string             `json:"trigger,omitempty"`
}

type InitSessionRequest struct {
	Repository               string `json:"repository,omitempty"`
	SelectedBranch           string `json:"selected_branch,omitempty"`
	InitialUserMsg           string `json:"initial_user_msg,omitempty"`
	ConversationInstructions string `json:"conversation_instructions,omitempty"`
	Trigger                  string `json:"trigger,omitempty"`
}

type ConversationResponse struct {
	Status             string             `json:"status"`
	ConversationID     string             `json:"conversation_id"`
	ConversationStatus ConversationStatus `json:"conversation_status,omitempty"`
	Message            string             `json:"message,omitempty"`
}

type ActionRequest struct {
	Action string `json:"action"`
	Args   string `json:"args"`
}
