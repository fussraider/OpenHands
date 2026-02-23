package handlers

import (
	"encoding/json"
	"net/http"
)

// MCPSSEHandler mocks the MCP SSE endpoint
func MCPSSEHandler(w http.ResponseWriter, r *http.Request) {
	// This should eventually be a WebSocket or SSE endpoint for FastMCP
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "MCP integration not implemented (mock)",
	})
}
