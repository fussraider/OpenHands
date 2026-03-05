package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"openhands-go/server/config"
	"openhands-go/server/models"
	"openhands-go/server/services"
	"openhands-go/server/store"
	"os"
	"testing"
)

func init() {
	config.AppConfig = &config.Config{
		AppMode: "oss",
		Sandbox: config.SandboxConfig{
			Runtime: "local",
		},
	}
	f, _ := os.CreateTemp("", "conversations.json")
	f.Close()
	ConversationStore = store.NewConversationStore(f.Name())

	f2, _ := os.CreateTemp("", "settings.json")
	f2.Close()
	SettingsStore = store.NewSettingsStore(f2.Name())

	RuntimeManager = services.NewRuntimeManager()
	// Mock broadcaster
	ActionService = services.NewActionService(ConversationStore, RuntimeManager, nil)
}

func TestNewConversationHandler(t *testing.T) {
	reqBody := models.InitSessionRequest{
		Repository: "test-repo",
	}
	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "/api/conversations", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(NewConversationHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var conversation models.ConversationInfo
	if err := json.NewDecoder(rr.Body).Decode(&conversation); err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}

	if conversation.SelectedRepository != "test-repo" {
		t.Errorf("handler returned unexpected repository: got %v want %v",
			conversation.SelectedRepository, "test-repo")
	}
}

func TestSearchConversationsHandler(t *testing.T) {
	// First create a conversation to ensure list is not empty
	reqBody := models.InitSessionRequest{Repository: "test-repo-2"}
	ConversationStore.CreateConversation(reqBody)

	req, err := http.NewRequest("GET", "/api/conversations", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SearchConversationsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}

	results, ok := response["results"].([]interface{})
	if !ok {
		t.Errorf("handler returned JSON missing 'results' array")
	}

	if len(results) == 0 {
		t.Errorf("handler returned empty list")
	}
}

func TestGetConversationHandler(t *testing.T) {
	// Create a conversation
	reqBody := models.InitSessionRequest{Repository: "test-repo-3"}
	created, _ := ConversationStore.CreateConversation(reqBody)

	// Create request with path value
	req, err := http.NewRequest("GET", "/api/conversations/"+created.ConversationID, nil)
	if err != nil {
		t.Fatal(err)
	}
	// Manually set path value for testing since httptest doesn't route
	req.SetPathValue("id", created.ConversationID)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetConversationHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var conversation models.ConversationInfo
	if err := json.NewDecoder(rr.Body).Decode(&conversation); err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}

	if conversation.ConversationID != created.ConversationID {
		t.Errorf("handler returned wrong conversation ID: got %v want %v",
			conversation.ConversationID, created.ConversationID)
	}
}
