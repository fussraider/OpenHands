package handlers

import (
	"encoding/json"
	"net/http"
)

func ConversationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Return empty list
	json.NewEncoder(w).Encode([]interface{}{})
}

func RepositoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]interface{}{})
}
