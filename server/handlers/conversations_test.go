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
	store.InitDB("file::memory:?cache=shared")
	ConversationStore = store.NewConversationStore()

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

func TestUpdateConversationHandler(t *testing.T) {
	// Create a conversation
	reqBody := models.InitSessionRequest{Repository: "test-repo-update"}
	created, _ := ConversationStore.CreateConversation(reqBody)

	// Update title
	updateReq := struct {
		Title string `json:"title"`
	}{
		Title: "New Title 123",
	}
	body, _ := json.Marshal(updateReq)

	req, err := http.NewRequest("PATCH", "/api/conversations/"+created.ConversationID, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("id", created.ConversationID)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateConversationHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify update in store
	updated, err := ConversationStore.GetConversation(created.ConversationID)
	if err != nil {
		t.Fatal(err)
	}

	if updated.Title != "New Title 123" {
		t.Errorf("expected title to be 'New Title 123', got '%s'", updated.Title)
	}
}

func TestStartStopConversationHandler(t *testing.T) {
	reqBody := models.InitSessionRequest{Repository: "test-repo-start-stop"}
	created, _ := ConversationStore.CreateConversation(reqBody)

	// Stop
	req, _ := http.NewRequest("POST", "/api/conversations/"+created.ConversationID+"/stop", nil)
	req.SetPathValue("id", created.ConversationID)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(StopConversationHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Stop returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Start
	req, _ = http.NewRequest("POST", "/api/conversations/"+created.ConversationID+"/start", nil)
	req.SetPathValue("id", created.ConversationID)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(StartConversationHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Start returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestAddMessageHandler(t *testing.T) {
	reqBody := models.InitSessionRequest{Repository: "test-repo-message"}
	created, _ := ConversationStore.CreateConversation(reqBody)

	msgReq := struct {
		Message string `json:"message"`
	}{
		Message: "Hello from test!",
	}
	body, _ := json.Marshal(msgReq)

	req, _ := http.NewRequest("POST", "/api/conversations/"+created.ConversationID+"/message", bytes.NewBuffer(body))
	req.SetPathValue("id", created.ConversationID)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddMessageHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("AddMessage returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify it was added to the event stream
	es := ActionService.GetEventStream(created.ConversationID)
	eventsList := es.GetEvents()
	found := false
	for _, ev := range eventsList {
		if reqData, ok := ev.Content.(models.ActionRequest); ok {
			if reqData.Action == "message" && reqData.Args["content"] == "Hello from test!" {
				found = true
				break
			}
		}
	}
	if !found {
		t.Errorf("Message was not found in event stream")
	}
}

func TestConversationAdditionalEndpoints(t *testing.T) {
	reqBody := models.InitSessionRequest{Repository: "test-repo"}
	created, _ := ConversationStore.CreateConversation(reqBody)
	id := created.ConversationID

	// Microagents
	req, _ := http.NewRequest("GET", "/api/conversations/"+id+"/microagents", nil)
	req.SetPathValue("id", id)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetConversationMicroagentsHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GetConversationMicroagentsHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Remember Prompt
	req, _ = http.NewRequest("GET", "/api/conversations/"+id+"/remember-prompt", nil)
	req.SetPathValue("id", id)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(GetRememberPromptHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GetRememberPromptHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// VSCode URL
	req, _ = http.NewRequest("GET", "/api/conversations/123/vscode-url", nil)
	req.SetPathValue("id", "123")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(GetVSCodeURLHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GetVSCodeURLHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Global VSCode URL endpoint used by V1 UI
	req, _ = http.NewRequest("GET", "/api/vscode/url", nil)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(GetGlobalVSCodeURLHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GetGlobalVSCodeURLHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Web Hosts
	req, _ = http.NewRequest("GET", "/api/conversations/123/web-hosts", nil)
	req.SetPathValue("id", "123")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(GetWebHostsHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GetWebHostsHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Config
	req, _ = http.NewRequest("GET", "/api/conversations/123/config", nil)
	req.SetPathValue("id", "123")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(GetConversationConfigHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GetConversationConfigHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Events
	req, _ = http.NewRequest("GET", "/api/conversations/123/events", nil)
	req.SetPathValue("id", "123")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(GetConversationEventsHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GetConversationEventsHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Exp Config
	req, _ = http.NewRequest("POST", "/api/conversations/123/exp-config", nil)
	req.SetPathValue("id", "123")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ExpConfigHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("ExpConfigHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestDeleteConversationHandler(t *testing.T) {
	// Create a conversation
	reqBody := models.InitSessionRequest{Repository: "test-repo-delete"}
	created, _ := ConversationStore.CreateConversation(reqBody)

	// Create request to delete
	req, err := http.NewRequest("DELETE", "/api/conversations/"+created.ConversationID, nil)
	if err != nil {
		t.Fatal(err)
	}
	// Manually set path value
	req.SetPathValue("id", created.ConversationID)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteConversationHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify conversation is actually deleted
	_, err = ConversationStore.GetConversation(created.ConversationID)
	if err == nil {
		t.Errorf("expected conversation to be deleted, but it was found")
	}

	// Test deleting non-existent conversation
	req, _ = http.NewRequest("DELETE", "/api/conversations/non-existent-id", nil)
	req.SetPathValue("id", "non-existent-id")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for non-existent id: got %v want %v", status, http.StatusNotFound)
	}
}
