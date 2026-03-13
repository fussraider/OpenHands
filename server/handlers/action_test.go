package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"openhands-go/server/config"
	"openhands-go/server/models"
	"testing"
)

func TestExecuteActionHandler(t *testing.T) {
	// Initialize config for testing (avoid nil pointer dereference)
	config.AppConfig = &config.Config{
		Sandbox: config.SandboxConfig{
			Runtime: "local",
		},
	}

	reqBody := models.ActionRequest{
		Action: "run",
		Args:   map[string]interface{}{"command": "echo 'hello world'"},
	}
	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "/api/conversations/123/action", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ExecuteActionHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}
}
