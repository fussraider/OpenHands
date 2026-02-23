package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"openhands-go/server/models"
	"testing"
)

func TestGetSettingsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/settings", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetSettingsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var settings models.Settings
	if err := json.NewDecoder(rr.Body).Decode(&settings); err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}

	// Default might be "en" or read from existing settings.json if left over
	if settings.Language == "" {
		t.Errorf("handler returned empty language")
	}
}

func TestStoreSettingsHandler(t *testing.T) {
	newSettings := models.Settings{
		Language: "es",
		Agent:    "CustomAgent",
	}
	body, _ := json.Marshal(newSettings)
	req, err := http.NewRequest("POST", "/api/settings", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(StoreSettingsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Verify settings were updated
	reqGet, _ := http.NewRequest("GET", "/api/settings", nil)
	rrGet := httptest.NewRecorder()
	handlerGet := http.HandlerFunc(GetSettingsHandler)
	handlerGet.ServeHTTP(rrGet, reqGet)

	var settings models.Settings
	json.NewDecoder(rrGet.Body).Decode(&settings)

	if settings.Language != "es" {
		t.Errorf("handler failed to update language: got %v want %v",
			settings.Language, "es")
	}
}
