package handlers

import (
	"encoding/json"
	"net/http"
	"openhands-go/server/models"
	"openhands-go/server/store"
)

var settingsStore = store.NewSettingsStore("settings.json")

func GetSettingsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settingsStore.Get())
}

func StoreSettingsHandler(w http.ResponseWriter, r *http.Request) {
	var newSettings models.Settings
	if err := json.NewDecoder(r.Body).Decode(&newSettings); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := settingsStore.Update(newSettings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Settings stored"})
}
