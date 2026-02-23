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
		Key   string `json:"key"`
		Value string `json:"value"`
	}{
		Key:   "API_KEY",
		Value: "12345",
	}
	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "/api/secrets", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(StoreSecretHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("StoreSecretHandler returned wrong status code: got %v want %v",
			status, http.StatusOK)
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

	var keys []string
	if err := json.NewDecoder(rr.Body).Decode(&keys); err != nil {
		t.Errorf("GetSecretsHandler returned invalid JSON: %v", err)
	}

	found := false
	for _, k := range keys {
		if k == "API_KEY" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("API_KEY not found in secrets: %v", keys)
	}

	// 3. Delete secret
	req, err = http.NewRequest("DELETE", "/api/secrets/API_KEY", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Manually set path value
	req.SetPathValue("key", "API_KEY")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(DeleteSecretHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("DeleteSecretHandler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Verify deletion
	secretsMutex.RLock()
	_, exists := secretsStore["API_KEY"]
	secretsMutex.RUnlock()

	if exists {
		t.Errorf("API_KEY was not deleted")
	}
}
