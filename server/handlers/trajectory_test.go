package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTrajectoryHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/conversations/123/trajectory", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTrajectoryHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}

	if _, ok := resp["trajectory"]; !ok {
		t.Errorf("response missing trajectory field")
	}

	// Should be empty list for now
	trajectory, ok := resp["trajectory"].([]interface{})
	if !ok {
		t.Errorf("trajectory is not a list")
	}
	if len(trajectory) != 0 {
		t.Errorf("trajectory is not empty")
	}
}
