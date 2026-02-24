package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"openhands-go/server/models"
	"openhands-go/server/services"
)

var (
	RuntimeManager *services.RuntimeManager
	ActionService  *services.ActionService
)

func ExecuteActionHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req models.ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := ActionService.ExecuteAction(r.Context(), id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"output": output})
}

// ProcessSocketAction handles actions coming from Socket.IO
func ProcessSocketAction(conversationID string, req models.ActionRequest) error {
	// Use background context as socket actions are async/long-lived
	ctx := context.Background()
	_, err := ActionService.ExecuteAction(ctx, conversationID, req)
	return err
}
