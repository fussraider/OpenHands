package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

	type customSecret struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	customSecrets := make([]customSecret, 0, len(secretsStore))
	for k := range secretsStore {
		// Filter out internal secrets like git_provider_*
		if !strings.HasPrefix(k, "git_provider_") {
			customSecrets = append(customSecrets, customSecret{
				Name:        k,
				Description: "", // Description not currently stored in simple MVP map
			})
		}
	}

	if customSecrets == nil {
		customSecrets = []customSecret{}
	}

	response := map[string]interface{}{
		"custom_secrets": customSecrets,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

	encryptedValue, err := encrypt(req.Value)
	if err != nil {
		http.Error(w, "Failed to encrypt secret", http.StatusInternalServerError)
		return
	}

	secretsMutex.Lock()
	defer secretsMutex.Unlock()
	secretsStore[req.Key] = encryptedValue

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
		encryptedValue, err := encrypt(string(data))
		if err == nil {
			secretsStore[key] = encryptedValue
		}
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

var encryptionKey []byte

func init() {
	key := os.Getenv("OPENHANDS_SECRETS_KEY")
	if key == "" {
		// Use a default key for local development if not provided,
		// but in production this should be set securely.
		// WARNING: This is highly insecure and should never be used in production.
		key = "default-insecure-32byte-secret!!"
	}

	// Ensure key is 32 bytes for AES-256
	if len(key) < 32 {
		padding := make([]byte, 32-len(key))
		for i := range padding {
			padding[i] = '0'
		}
		key += string(padding)
	} else if len(key) > 32 {
		key = key[:32]
	}

	encryptionKey = []byte(key)
}

func encrypt(text string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(text))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func decrypt(cryptoText string) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func GetSecret(key string) (string, bool) {
	secretsMutex.RLock()
	defer secretsMutex.RUnlock()

	encryptedValue, ok := secretsStore[key]
	if !ok {
		return "", false
	}

	decryptedValue, err := decrypt(encryptedValue)
	if err != nil {
		// If decryption fails, it might be an old unencrypted secret or a bad key.
		// For safety, return the encrypted string or empty.
		return "", false
	}
	return decryptedValue, true
}
