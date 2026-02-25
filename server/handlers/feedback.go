package handlers

import (
	"encoding/json"
	"net/http"
)

type FeedbackRequest struct {
	Email       string `json:"email"`
	Version     string `json:"version"`
	Permissions string `json:"permissions"`
	Polarity    string `json:"polarity"`
	Feedback    string `json:"feedback"`
}

type FeedbackResponse struct {
	Status string `json:"status"`
}

func SubmitFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	// conversationID := r.PathValue("id")

	var req FeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Mock storing feedback:
	// In reality, this would store the feedback and the trajectory (events).
	// For now, we just return success.

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(FeedbackResponse{Status: "ok"})
}
