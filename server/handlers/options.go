package handlers

import (
	"encoding/json"
	"net/http"
)

import (
	"openhands-go/server/agent"
	"openhands-go/server/security"
)

func ModelsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// In Python, this fetches dynamically from LiteLLM.
	// For Go, without a direct LiteLLM registry, we provide a robust static list.
	models := []string{
		"gpt-4o", "gpt-4-turbo", "gpt-3.5-turbo",
		"claude-3-5-sonnet-20240620", "claude-3-opus-20240229",
		"gemini-1.5-pro", "gemini-1.5-flash",
		"llama3", "mistral",
	}
	json.NewEncoder(w).Encode(models)
}

func GetAgentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// For Go, we expose the agents we have implemented in the agent package
	agents := agent.ListAgents()
	json.NewEncoder(w).Encode(agents)
}

func GetSecurityAnalyzersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	analyzers := security.ListAnalyzers()
	json.NewEncoder(w).Encode(analyzers)
}
