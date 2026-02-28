package handlers

import (
	"encoding/json"
	"net/http"
)

// SecurityAPIHandler simulates the Security API endpoint
func SecurityAPIHandler(w http.ResponseWriter, r *http.Request) {
	// Revert to 501 Not Implemented because there is no frontend/backend
	// implementation of a user-facing security approval API in Go yet.
	// BasicAnalyzer and LLMSecurityAnalyzer are integrated directly into the Agent loop.
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "Security analyzer not initialized (mock)",
	})
}
