package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1SandboxRoutes(t *testing.T) {
	mux := http.NewServeMux()
	RegisterV1Routes(mux)

	tests := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/sandboxes/search"},
		{"GET", "/api/v1/sandboxes"},
		{"POST", "/api/v1/sandboxes"},
		{"POST", "/api/v1/sandboxes/123/pause"},
		{"POST", "/api/v1/sandboxes/123/resume"},
		{"DELETE", "/api/v1/sandboxes/123"},
	}

	for _, tt := range tests {
		req, _ := http.NewRequest(tt.method, tt.path, nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("%s %s returned wrong status code: got %v want %v", tt.method, tt.path, rr.Code, http.StatusOK)
		}

		var resp interface{}
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Errorf("%s %s returned invalid JSON: %v", tt.method, tt.path, err)
		}
	}
}
