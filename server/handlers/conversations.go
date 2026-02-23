package handlers

import (
	"encoding/json"
	"net/http"
	"openhands-go/server/models"
	"openhands-go/server/store"
)

// Global instance for now
var globalConversationStore = store.NewConversationStore()

func SearchConversationsHandler(w http.ResponseWriter, r *http.Request) {
	conversations := globalConversationStore.ListConversations()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)
}

func NewConversationHandler(w http.ResponseWriter, r *http.Request) {
	var req models.InitSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conversation, err := globalConversationStore.CreateConversation(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversation)
}

func GetConversationHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	conversation, err := globalConversationStore.GetConversation(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversation)
}
