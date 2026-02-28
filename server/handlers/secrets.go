package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
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

// StoreGitProvidersHandler handles POST /api/add-git-providers
func StoreGitProvidersHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProviderTokens map[string]struct {
			Token  string `json:"token"`
			UserID string `json:"user_id"`
			Host   string `json:"host"`
		} `json:"provider_tokens"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	secretsMutex.Lock()
	defer secretsMutex.Unlock()

	// Store these tokens in the secretsStore with a prefix
	for provider, tokenData := range req.ProviderTokens {
		key := "git_provider_" + provider
		// Simplification for MVP: just store the token value.
		// Real implementation might store JSON string representing the struct.
		data, _ := json.Marshal(tokenData)
		secretsStore[key] = string(data)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Git providers stored"})
}

// UnsetGitProvidersHandler handles POST /api/unset-provider-tokens
func UnsetGitProvidersHandler(w http.ResponseWriter, r *http.Request) {
	secretsMutex.Lock()
	defer secretsMutex.Unlock()

	// Find and delete all git provider tokens
	for key := range secretsStore {
		if strings.HasPrefix(key, "git_provider_") {
			delete(secretsStore, key)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Unset Git provider tokens"})
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
