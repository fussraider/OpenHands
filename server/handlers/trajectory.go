package handlers

import (
	"encoding/json"
	"net/http"
)

type Trajectory struct {
	Events []interface{} `json:"events"`
}

func GetTrajectoryHandler(w http.ResponseWriter, r *http.Request) {
	// conversationID := r.PathValue("id")

	// Mock: return empty trajectory
	// In the real system, this would retrieve events from the event store.
	trajectory := Trajectory{
		Events: []interface{}{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"trajectory": trajectory.Events})
}
