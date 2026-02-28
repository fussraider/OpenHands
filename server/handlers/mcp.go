package handlers

import (
	"encoding/json"
	"net/http"
)

// MCPSSEHandler mocks the MCP SSE endpoint
func MCPSSEHandler(w http.ResponseWriter, r *http.Request) {
	// Revert to 501 Not Implemented because the SSE route is currently unneeded for native
	// Go Stdio MCP implementation, but the API endpoint itself is required for parity checking.
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "MCP integration not implemented (mock)",
	})
}
