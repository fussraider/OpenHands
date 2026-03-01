package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"openhands-go/server/models"
	"openhands-go/server/store"
	"os"
	"testing"
)

func TestGetMicroagentManagementConversationsHandler(t *testing.T) {
	// Setup test conversation store
	f, _ := os.CreateTemp("", "conversations_test.jsonl")
	defer os.Remove(f.Name())
	ConversationStore = store.NewConversationStore(f.Name())

	// Create some conversations
	ConversationStore.CreateConversation(models.InitSessionRequest{
		Repository: "owner/repo1",
		Trigger:    "MICROAGENT_MANAGEMENT",
	})
	ConversationStore.CreateConversation(models.InitSessionRequest{
		Repository: "owner/repo2",
		Trigger:    "MICROAGENT_MANAGEMENT",
	})
	ConversationStore.CreateConversation(models.InitSessionRequest{
		Repository: "owner/repo1",
		Trigger:    "other_trigger",
	})

	// Test 1: Get all microagent conversations
	req, _ := http.NewRequest("GET", "/api/microagent-management/conversations", nil)
	rr := httptest.NewRecorder()
	http.HandlerFunc(GetMicroagentManagementConversationsHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&response)

	results := response["results"].([]interface{})
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

	// Test 2: Filter by repository
	req, _ = http.NewRequest("GET", "/api/microagent-management/conversations?selected_repository=owner/repo1", nil)
	rr = httptest.NewRecorder()
	http.HandlerFunc(GetMicroagentManagementConversationsHandler).ServeHTTP(rr, req)

	json.NewDecoder(rr.Body).Decode(&response)
	results = response["results"].([]interface{})
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}
