package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"openhands-go/server/models"
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

func TestGetWebClientConfigHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/web-client/config", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetWebClientConfigHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var config models.WebClientConfig
	err = json.Unmarshal(rr.Body.Bytes(), &config)
	if err != nil {
		t.Errorf("could not unmarshal response body: %v", err)
	}

	if config.AppMode != "oss" {
		t.Errorf("expected AppMode oss, got %s", config.AppMode)
	}
}
