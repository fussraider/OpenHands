package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"openhands-go/server/config"
	"os"
	"path/filepath"
	"sync"
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

var feedbackMu sync.Mutex

func SubmitFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	var req FeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	slog.Debug("Got feedback", "email", req.Email, "polarity", req.Polarity, "permissions", req.Permissions)

	// Store feedback in file
	if config.AppConfig.FileStorePath != "" {
		feedbackPath := filepath.Join(config.AppConfig.FileStorePath, "feedback.jsonl")

		feedbackMu.Lock()
		defer feedbackMu.Unlock()

		f, err := os.OpenFile(feedbackPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			defer f.Close()
			data, _ := json.Marshal(req)
			f.Write(data)
			f.Write([]byte("\n"))
			slog.Debug("Stored feedback", "path", feedbackPath)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(FeedbackResponse{Status: "ok"})
}
