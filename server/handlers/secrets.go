package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Secrets map[string]string

var (
	secretsStore = make(Secrets)
	secretsMutex sync.RWMutex
)

// GetSecretsHandler lists secrets (keys only, masked values?)
// Python implementation returns all secrets, masked if exposed
func GetSecretsHandler(w http.ResponseWriter, r *http.Request) {
	secretsMutex.RLock()
	defer secretsMutex.RUnlock()

	// In legacy V0, /api/secrets returns a list of secret names
	keys := make([]string, 0, len(secretsStore))
	for k := range secretsStore {
		keys = append(keys, k)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keys)
}

// StoreSecretHandler stores a secret
func StoreSecretHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	secretsMutex.Lock()
	defer secretsMutex.Unlock()
	secretsStore[req.Key] = req.Value

	w.WriteHeader(http.StatusOK)
}

// DeleteSecretHandler deletes a secret
func DeleteSecretHandler(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		return
	}

	secretsMutex.Lock()
	defer secretsMutex.Unlock()
	delete(secretsStore, key)

	w.WriteHeader(http.StatusOK)
}
