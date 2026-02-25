package store

import (
	"encoding/json"
	"openhands-go/server/models"
	"os"
	"sync"
)

type SettingsStore struct {
	mu       sync.RWMutex
	settings models.Settings
	filePath string
}

func NewSettingsStore(filePath string) *SettingsStore {
	store := &SettingsStore{
		settings: models.DefaultSettings(),
		filePath: filePath,
	}
	store.load()
	return store
}

func (s *SettingsStore) load() {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return
	}
	json.Unmarshal(data, &s.settings)
}

func (s *SettingsStore) save() error {
	data, err := json.MarshalIndent(s.settings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0644)
}

func (s *SettingsStore) Get() models.Settings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.settings
}

func (s *SettingsStore) Update(newSettings models.Settings) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.settings = newSettings
	return s.save()
}
