package handlers

import (
	"encoding/json"
	"net/http"
	"openhands-go/server/models"
	"openhands-go/server/services"
)

var (
	runtimeManager = services.NewRuntimeManager()
	actionService  = services.NewActionService(conversationStore, runtimeManager)
)

func ExecuteActionHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req models.ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := actionService.ExecuteAction(r.Context(), id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"output": output})
}
