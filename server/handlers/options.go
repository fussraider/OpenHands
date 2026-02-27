package handlers

import (
	"encoding/json"
	"net/http"
)

func ModelsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Return a dummy model list
	models := []string{"gpt-4", "gpt-3.5-turbo", "claude-3-opus"}
	json.NewEncoder(w).Encode(models)
}

func GetAgentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Return mock agents list
	agents := []string{"CodeActAgent", "MonologueAgent", "PlannerAgent"}
	json.NewEncoder(w).Encode(agents)
}

func GetSecurityAnalyzersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Return mock security analyzers
	analyzers := []string{"mitmproxy", "llm", "invariant"}
	json.NewEncoder(w).Encode(analyzers)
}
