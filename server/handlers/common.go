package handlers

import (
	"encoding/json"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func ModelsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Return a dummy model list
	models := []string{"gpt-4", "gpt-3.5-turbo", "claude-3-opus"}
	json.NewEncoder(w).Encode(models)
}
