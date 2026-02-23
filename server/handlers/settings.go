package handlers

import (
	"encoding/json"
	"net/http"
	"openhands-go/server/models"
	"sync"
)

var (
	settingsStore = models.DefaultSettings()
	settingsMutex sync.RWMutex
)

func GetSettingsHandler(w http.ResponseWriter, r *http.Request) {
	settingsMutex.RLock()
	defer settingsMutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settingsStore)
}

func StoreSettingsHandler(w http.ResponseWriter, r *http.Request) {
	var newSettings models.Settings
	if err := json.NewDecoder(r.Body).Decode(&newSettings); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	// In a real implementation, we would merge with existing settings
	// For now, we just overwrite
	settingsStore = newSettings

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Settings stored"})
}
