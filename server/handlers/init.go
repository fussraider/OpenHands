package handlers

import (
	"openhands-go/server/config"
	"openhands-go/server/services"
	"openhands-go/server/store"
	"path/filepath"
)

func InitHandlers() {
	settingsPath := "settings.json"
	dbPath := "openhands.db"
	if config.AppConfig.FileStorePath != "" {
		settingsPath = filepath.Join(config.AppConfig.FileStorePath, "settings.json")
		dbPath = filepath.Join(config.AppConfig.FileStorePath, "openhands.db")
	}
	if err := store.InitDB(dbPath); err != nil {
		panic(err)
	}

	ConversationStore = store.NewConversationStore()
	SettingsStore = store.NewSettingsStore(settingsPath)

	userSettings := SettingsStore.Get()
	if userSettings.LLMAPIKey != "" {
		config.AppConfig.LLM.APIKey = userSettings.LLMAPIKey
	}
	if userSettings.LLMModel != "" {
		config.AppConfig.LLM.Model = userSettings.LLMModel
	}
	if userSettings.LLMBaseURL != "" {
		config.AppConfig.LLM.BaseURL = userSettings.LLMBaseURL
	}

	RuntimeManager = services.NewRuntimeManager()
	// No broadcaster needed — V1 WebSocket clients subscribe to EventStream directly
	ActionService = services.NewActionService(ConversationStore, RuntimeManager, nil)
	GitService = services.NewGitService(config.AppConfig)
}
