package handlers

import (
	"openhands-go/server/config"
	"openhands-go/server/services"
	"openhands-go/server/store"
	"openhands-go/server/ws"
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
	ActionService = services.NewActionService(ConversationStore, RuntimeManager, ws.BroadcastEvent)
	GitService = services.NewGitService(config.AppConfig)
}
