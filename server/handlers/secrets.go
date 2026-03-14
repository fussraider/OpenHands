package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"openhands-go/server/models"
	"os"
	"strings"
)

// GetSecretsHandler lists secrets (keys only, masked values?)
// Python implementation returns all secrets, masked if exposed
func GetSecretsHandler(w http.ResponseWriter, r *http.Request) {
	allSecrets, err := SecretsStore.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type customSecret struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	customSecrets := make([]customSecret, 0, len(allSecrets))
	for _, sec := range allSecrets {
		// Filter out internal secrets like git_provider_*
		if !strings.HasPrefix(sec.Name, "git_provider_") {
			customSecrets = append(customSecrets, customSecret{
				Name:        sec.Name,
				Description: sec.Description,
			})
		}
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
		Name        string `json:"name"`
		Value       string `json:"value"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "secret name is required", http.StatusBadRequest)
		return
	}

	encryptedValue, err := encrypt(req.Value)
	if err != nil {
		http.Error(w, "Failed to encrypt secret", http.StatusInternalServerError)
		return
	}

	err = SecretsStore.Save(&models.SecretInfo{
		Name:        req.Name,
		Value:       encryptedValue,
		Description: req.Description,
	})

	if err != nil {
		http.Error(w, "Failed to store secret", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
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

	// Store these tokens in the secretsStore with a prefix
	for provider, tokenData := range req.ProviderTokens {
		key := "git_provider_" + provider
		// Simplification for MVP: just store the token value.
		// Real implementation might store JSON string representing the struct.
		data, _ := json.Marshal(tokenData)
		encryptedValue, err := encrypt(string(data))
		if err == nil {
			SecretsStore.Save(&models.SecretInfo{
				Name:  key,
				Value: encryptedValue,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Git providers stored"})
}

// UnsetGitProvidersHandler handles POST /api/unset-provider-tokens
func UnsetGitProvidersHandler(w http.ResponseWriter, r *http.Request) {
	allSecrets, err := SecretsStore.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Find and delete all git provider tokens
	for _, sec := range allSecrets {
		if strings.HasPrefix(sec.Name, "git_provider_") {
			SecretsStore.Delete(sec.Name)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Unset Git provider tokens"})
}

// UpdateSecretHandler updates a secret
func UpdateSecretHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "secret id is required", http.StatusBadRequest)
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "secret name is required", http.StatusBadRequest)
		return
	}

	existingSecret, err := SecretsStore.Get(id)
	if err != nil {
		http.Error(w, "secret not found", http.StatusNotFound)
		return
	}

	// If the name changed, we need to move the key.
	if id != req.Name {
		// Check if new name already exists
		if _, err := SecretsStore.Get(req.Name); err == nil {
			http.Error(w, "a secret with the new name already exists", http.StatusBadRequest)
			return
		}

		// Save under new name
		err = SecretsStore.Save(&models.SecretInfo{
			Name:        req.Name,
			Value:       existingSecret.Value,
			Description: req.Description,
		})
		if err == nil {
			SecretsStore.Delete(id)
		}
	} else {
		// Update description
		existingSecret.Description = req.Description
		SecretsStore.Save(existingSecret)
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteSecretHandler deletes a secret
func DeleteSecretHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "secret id is required", http.StatusBadRequest)
		return
	}

	if _, err := SecretsStore.Get(id); err != nil {
		http.Error(w, "secret not found", http.StatusNotFound)
		return
	}

	SecretsStore.Delete(id)

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
	if SecretsStore == nil {
		return "", false
	}
	secretInfo, err := SecretsStore.Get(key)
	if err != nil {
		return "", false
	}

	decryptedValue, err := decrypt(secretInfo.Value)
	if err != nil {
		// If decryption fails, it might be an old unencrypted secret or a bad key.
		// For safety, return the encrypted string or empty.
		return "", false
	}

	slog.Debug("Loaded token from secret store", "secret_name", key)

	return decryptedValue, true
}
