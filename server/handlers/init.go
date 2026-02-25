package handlers

import (
	"openhands-go/server/config"
	"openhands-go/server/services"
	"openhands-go/server/store"
	"openhands-go/server/ws"
	"path/filepath"
)

func InitHandlers() {
	convPath := "conversations.json"
	settingsPath := "settings.json"
	if config.AppConfig.FileStorePath != "" {
		convPath = filepath.Join(config.AppConfig.FileStorePath, "conversations.json")
		settingsPath = filepath.Join(config.AppConfig.FileStorePath, "settings.json")
	}
	ConversationStore = store.NewConversationStore(convPath)
	SettingsStore = store.NewSettingsStore(settingsPath)

	RuntimeManager = services.NewRuntimeManager()
	ActionService = services.NewActionService(ConversationStore, RuntimeManager, ws.BroadcastEvent)
	GithubService = services.NewGithubService(config.AppConfig.Github)
}
