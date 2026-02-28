package handlers

import (
	"encoding/json"
	"net/http"
)

func ModelsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Return a static list of commonly supported models
	models := []string{"gpt-4", "gpt-3.5-turbo", "claude-3-opus", "claude-3-sonnet"}
	json.NewEncoder(w).Encode(models)
}

func GetAgentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Return dynamically supported agents based on Go implementation
	agents := []string{"CodeActAgent", "BrowsingAgent"}
	json.NewEncoder(w).Encode(agents)
}

func GetSecurityAnalyzersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Return dynamically supported security analyzers based on Go implementation
	analyzers := []string{"basic", "llm"}
	json.NewEncoder(w).Encode(analyzers)
}
