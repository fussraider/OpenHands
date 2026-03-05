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

	response := map[string]interface{}{
		"results":      conversations,
		"next_page_id": nil,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

func StartConversationHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "conversation id required", http.StatusBadRequest)
		return
	}

	conversation, err := ConversationStore.GetConversation(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	ctx := context.Background()

	// Ensure runtime is created if not exists
	_, err = RuntimeManager.GetRuntime(id)
	if err != nil {
		_, err = RuntimeManager.CreateRuntime(ctx, id)
		if err != nil {
			http.Error(w, "Failed to create runtime: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Ensure agent loop is started
	es := ActionService.GetEventStream(id)
	err = RuntimeManager.StartAgent(ctx, id, es)
	if err != nil {
		http.Error(w, "Failed to start agent: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversation)
}

func StopConversationHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "conversation id required", http.StatusBadRequest)
		return
	}

	conversation, err := ConversationStore.GetConversation(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if RuntimeManager != nil {
		err = RuntimeManager.StopRuntime(id)
		if err != nil {
			http.Error(w, "Failed to stop runtime: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversation)
}

func UpdateConversationHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "conversation id required", http.StatusBadRequest)
		return
	}

	var req struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	err := ConversationStore.UpdateConversation(id, req.Title)
	if err != nil {
		if err.Error() == "conversation not found" {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(true)
}

func DeleteConversationHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "conversation id required", http.StatusBadRequest)
		return
	}

	err := ConversationStore.DeleteConversation(id)
	if err != nil {
		if err.Error() == "conversation not found" {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Also stop the agent loop / runtime if it exists
	if RuntimeManager != nil {
		RuntimeManager.StopRuntime(id)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
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
