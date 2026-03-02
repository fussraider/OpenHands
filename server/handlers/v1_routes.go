package handlers

import (
	"encoding/json"
	"net/http"
	"openhands-go/server/services"
)

// The v1 API in Python (app_server) handles advanced enterprise-like features
// such as sandboxes, user auth context, web client events, and persistent app_conversations.
// For the open-source Go backend, we will implement the foundational endpoints
// required for frontend compatibility, connecting them to existing local systems.

var runtimeManager *services.RuntimeManager

func SetRuntimeManager(rm *services.RuntimeManager) {
	runtimeManager = rm
}

func RegisterV1Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/sandboxes/search", V1SearchSandboxesHandler)
	mux.HandleFunc("GET /api/v1/sandboxes", V1BatchGetSandboxesHandler)
	mux.HandleFunc("POST /api/v1/sandboxes", V1StartSandboxHandler)
	mux.HandleFunc("POST /api/v1/sandboxes/{id}/pause", V1PauseSandboxHandler)
	mux.HandleFunc("POST /api/v1/sandboxes/{id}/resume", V1ResumeSandboxHandler)
	mux.HandleFunc("DELETE /api/v1/sandboxes/{id}", V1DeleteSandboxHandler)

	// Events
	mux.HandleFunc("GET /api/v1/conversation/{conversation_id}/events/search", V1SearchEventsHandler)
	mux.HandleFunc("GET /api/v1/conversation/{conversation_id}/events", V1BatchGetEventsHandler)
	mux.HandleFunc("GET /api/v1/conversation/{conversation_id}/events/count", V1CountEventsHandler)
}

func V1SearchSandboxesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"items": []interface{}{},
		"next_page_id": nil,
	})
}

func V1BatchGetSandboxesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]interface{}{})
}

func V1StartSandboxHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure a runtime is created for this session
	sandboxID := "local-sandbox"

	if runtimeManager != nil {
		// For true V1 parity, the payload contains config options for the sandbox
		// Since this is HTTP context, we use the request context to bind the lifecycle
		ctx := r.Context()

		// Attempt to create runtime. If it fails, return 500
		_, err := runtimeManager.CreateRuntime(ctx, sandboxID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id": sandboxID,
		"status": "RUNNING",
	})
}

func V1PauseSandboxHandler(w http.ResponseWriter, r *http.Request) {
	// MVP implementation: Pausing would effectively stop the agent loop.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func V1ResumeSandboxHandler(w http.ResponseWriter, r *http.Request) {
	// MVP implementation: Resume the agent loop for the corresponding session.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func V1DeleteSandboxHandler(w http.ResponseWriter, r *http.Request) {
	sandboxID := r.PathValue("id")
	if runtimeManager != nil && sandboxID != "" {
		runtimeManager.StopRuntime(sandboxID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func V1SearchEventsHandler(w http.ResponseWriter, r *http.Request) {
	conversationID := r.PathValue("conversation_id")
	if conversationID == "" {
		http.Error(w, "conversation_id required", http.StatusBadRequest)
		return
	}

	// This is a minimal bridge to the EventStream for v1 parity
	// In Python this hits event_service.search_events which reads from SQLite/DB
	// For local open-source we just use the in-memory or file-backed EventStream

	w.Header().Set("Content-Type", "application/json")

	// Mock empty page format expected by frontend
	json.NewEncoder(w).Encode(map[string]interface{}{
		"items": []interface{}{},
		"next_page_id": nil,
	})
}

func V1BatchGetEventsHandler(w http.ResponseWriter, r *http.Request) {
	conversationID := r.PathValue("conversation_id")
	if conversationID == "" {
		http.Error(w, "conversation_id required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]interface{}{})
}

func V1CountEventsHandler(w http.ResponseWriter, r *http.Request) {
	conversationID := r.PathValue("conversation_id")
	if conversationID == "" {
		http.Error(w, "conversation_id required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(0)
}
