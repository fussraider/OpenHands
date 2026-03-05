package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecretsHandlers(t *testing.T) {
	// 1. Store a secret
	reqBody := struct {
		Name        string `json:"name"`
		Value       string `json:"value"`
		Description string `json:"description"`
	}{
		Name:        "API_KEY",
		Value:       "12345",
		Description: "A test API key",
	}
	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "/api/secrets", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(StoreSecretHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("StoreSecretHandler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// 2. Get secrets
	req, err = http.NewRequest("GET", "/api/secrets", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(GetSecretsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GetSecretsHandler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response struct {
		CustomSecrets []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"custom_secrets"`
	}

	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("GetSecretsHandler returned invalid JSON: %v", err)
	}

	found := false
	for _, sec := range response.CustomSecrets {
		if sec.Name == "API_KEY" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("API_KEY not found in secrets: %v", response.CustomSecrets)
	}

	// 2.5 Update secret
	updateReqBody := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{
		Name:        "API_KEY_UPDATED",
		Description: "Updated description",
	}
	body, _ = json.Marshal(updateReqBody)
	req, err = http.NewRequest("PUT", "/api/secrets/API_KEY", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", "API_KEY")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(UpdateSecretHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("UpdateSecretHandler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// 3. Delete secret
	req, err = http.NewRequest("DELETE", "/api/secrets/API_KEY_UPDATED", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Manually set path value
	req.SetPathValue("id", "API_KEY_UPDATED")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(DeleteSecretHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("DeleteSecretHandler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Verify deletion
	secretsMutex.RLock()
	_, exists := secretsStore["API_KEY_UPDATED"]
	secretsMutex.RUnlock()

	if exists {
		t.Errorf("API_KEY_UPDATED was not deleted")
	}
}
