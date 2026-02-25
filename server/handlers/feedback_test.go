package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubmitFeedbackHandler(t *testing.T) {
	reqBody := FeedbackRequest{
		Email:    "test@example.com",
		Polarity: "positive",
		Feedback: "Great tool!",
	}
	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "/api/conversations/123/submit-feedback", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SubmitFeedbackHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp FeedbackResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}

	if resp.Status != "ok" {
		t.Errorf("handler returned unexpected status: got %v want %v",
			resp.Status, "ok")
	}
}
