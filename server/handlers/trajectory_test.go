package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"openhands-go/server/events"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGetTrajectoryHandler(t *testing.T) {
	// Initialize handler dependencies (via init() in conversations_test.go or similar)
	// We assume init() runs.

	// Create a conversation and add events
	convID := "traj-test-conv"
	es := ActionService.GetEventStream(convID)
	es.AddEvent(events.Event{
		ID: uuid.New().String(),
		Type: events.EventTypeAction,
		Content: "test",
		Timestamp: time.Now(),
	})

	req, _ := http.NewRequest("GET", "/api/conversations/"+convID+"/trajectory", nil)
	req.SetPathValue("id", convID)

	rr := httptest.NewRecorder()
	http.HandlerFunc(GetTrajectoryHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	var resp map[string][]events.Event
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Errorf("Handler returned invalid JSON: %v", err)
	}

	if len(resp["trajectory"]) != 1 {
		t.Errorf("Expected 1 event, got %d", len(resp["trajectory"]))
	}
}
