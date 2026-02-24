package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"openhands-go/server/models"
	"openhands-go/server/store"
)

var conversationStore = store.NewConversationStore("conversations.json")

func SearchConversationsHandler(w http.ResponseWriter, r *http.Request) {
	conversations := conversationStore.ListConversations()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)
}

func NewConversationHandler(w http.ResponseWriter, r *http.Request) {
	var req models.InitSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conversation, err := conversationStore.CreateConversation(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Initialize Runtime and Agent
	// Use background context for runtime/agent lifecycle, as they outlive the request
	ctx := context.Background()
	_, err = runtimeManager.CreateRuntime(ctx, conversation.ConversationID)
	if err != nil {
		// Log error but don't fail conversation creation?
		// Or fail.
		http.Error(w, "Failed to create runtime: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Start Agent Loop
	// Need access to EventStream from ActionService?
	es := actionService.GetEventStream(conversation.ConversationID)
	err = runtimeManager.StartAgent(ctx, conversation.ConversationID, es)
	if err != nil {
		http.Error(w, "Failed to start agent: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversation)
}

func GetConversationHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	conversation, err := conversationStore.GetConversation(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversation)
}
