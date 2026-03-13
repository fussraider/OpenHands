package handlers

import (
	"encoding/json"
	"net/http"
	"openhands-go/server/config"
	"openhands-go/server/models"
	"openhands-go/server/store"
)

var SettingsStore *store.SettingsStore

func GetSettingsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SettingsStore.Get())
}

func StoreSettingsHandler(w http.ResponseWriter, r *http.Request) {
	var newSettings models.Settings
	if err := json.NewDecoder(r.Body).Decode(&newSettings); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := SettingsStore.Update(newSettings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if newSettings.LLMAPIKey != "" {
		config.AppConfig.LLM.APIKey = newSettings.LLMAPIKey
	}
	if newSettings.LLMModel != "" {
		config.AppConfig.LLM.Model = newSettings.LLMModel
	}
	if newSettings.LLMBaseURL != "" {
		config.AppConfig.LLM.BaseURL = newSettings.LLMBaseURL
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Settings stored"})
}
