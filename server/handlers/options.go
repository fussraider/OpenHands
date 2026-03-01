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

	// Pre-defined set of openhands models representing what is usually supported by proxies like litellm
	openhandsModels := []string{
		"openhands/claude-opus-4-5-20251101",
		"openhands/claude-sonnet-4-5-20250929",
		"openhands/gpt-5.2-codex",
		"openhands/gpt-5.2",
		"openhands/minimax-m2.5",
		"openhands/gemini-3-pro-preview",
		"openhands/gemini-3-flash-preview",
		"openhands/deepseek-chat",
		"openhands/devstral-medium-2512",
		"openhands/kimi-k2-0711-preview",
		"openhands/qwen3-coder-480b",
		"claude-3-5-sonnet-20241022",
		"gpt-4o",
	}

	models := make([]string, 0, len(openhandsModels))
	models = append(models, openhandsModels...)

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
