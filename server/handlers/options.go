package handlers

import (
	"encoding/json"
	"net/http"
	"openhands-go/server/config"
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
	// For Go implementation, these are the natively defined agents.
	// Since there is no dynamic registry in Go yet (they are hardcoded in factories),
	// we explicitly list the ones compiled into the binary.
	agents := []string{"CodeActAgent", "BrowsingAgent"}
	json.NewEncoder(w).Encode(agents)
}

func GetSecurityAnalyzersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Explicitly list the compiled-in security analyzers from server/security/
	analyzers := []string{"basic", "llm"}
	json.NewEncoder(w).Encode(analyzers)
}
