package handlers

import (
	"encoding/json"
	"net/http"
)

type Trajectory struct {
	Events []interface{} `json:"events"`
}

func GetTrajectoryHandler(w http.ResponseWriter, r *http.Request) {
	conversationID := r.PathValue("id")

	// Retrieve events from ActionService (which gets from EventStream)
	events := ActionService.GetEventStream(conversationID).GetEvents()

	// Convert to interface slice as expected by frontend
	// Note: Event struct matches expected JSON format mostly.
	// Frontend expects list of events.

	w.Header().Set("Content-Type", "application/json")
	// Wrapping in "trajectory" key as per Python implementation?
	// Python: `return {"trajectory": [event_to_dict(e) for e in events]}`
	json.NewEncoder(w).Encode(map[string]interface{}{"trajectory": events})
}
