package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"openhands-go/server/models"
	"openhands-go/server/store"
)

var ConversationStore *store.ConversationStore

func SearchConversationsHandler(w http.ResponseWriter, r *http.Request) {
	conversations := ConversationStore.ListConversations()
	if conversations == nil {
		conversations = []models.ConversationInfo{} // Ensure we return [] instead of null
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)
}

func NewConversationHandler(w http.ResponseWriter, r *http.Request) {
	var req models.InitSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conversation, err := ConversationStore.CreateConversation(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Initialize Runtime and Agent
	// Use background context for runtime/agent lifecycle, as they outlive the request
	ctx := context.Background()
	_, err = RuntimeManager.CreateRuntime(ctx, conversation.ConversationID)
	if err != nil {
		// Log error but don't fail conversation creation?
		// Or fail.
		http.Error(w, "Failed to create runtime: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Start Agent Loop
	// Need access to EventStream from ActionService?
	es := ActionService.GetEventStream(conversation.ConversationID)
	err = RuntimeManager.StartAgent(ctx, conversation.ConversationID, es)
	if err != nil {
		http.Error(w, "Failed to start agent: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversation)
}

func GetConversationHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	conversation, err := ConversationStore.GetConversation(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversation)
}
