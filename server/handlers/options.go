package handlers

import (
	"encoding/json"
	"net/http"
	"openhands-go/server/agent"
	"openhands-go/server/config"
	"openhands-go/server/security"
)

func ModelsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// If a specific provider is loaded, return its models. For langchaingo OpenAI we can't easily list,
	// so we return the currently configured model as guaranteed, plus common defaults if using standard providers.
	models := []string{"gpt-4", "gpt-3.5-turbo", "claude-3-opus", "claude-3-sonnet"}

	if config.AppConfig != nil && config.AppConfig.LLM.Model != "" {
		found := false
		for _, m := range models {
			if m == config.AppConfig.LLM.Model {
				found = true
				break
			}
		}
		if !found {
			models = append([]string{config.AppConfig.LLM.Model}, models...)
		}
	}

	json.NewEncoder(w).Encode(models)
}

func GetAgentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Fetch agents dynamically from the agent package registry.
	agents := agent.GetAvailableAgents()
	json.NewEncoder(w).Encode(agents)
}

func GetSecurityAnalyzersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Fetch analyzers dynamically from the security package registry.
	analyzers := security.GetAvailableAnalyzers()
	json.NewEncoder(w).Encode(analyzers)
}
