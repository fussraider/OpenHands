package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOptionsHandlers(t *testing.T) {
	// Test ModelsHandler
	req, _ := http.NewRequest("GET", "/api/options/models", nil)
	rr := httptest.NewRecorder()
	http.HandlerFunc(ModelsHandler).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("ModelsHandler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
	var models []string
	if err := json.NewDecoder(rr.Body).Decode(&models); err != nil {
		t.Errorf("ModelsHandler returned invalid JSON: %v", err)
	}
	if len(models) != 4 || models[0] != "gpt-4" {
		t.Errorf("ModelsHandler returned unexpected list: %v", models)
	}

	// Test GetAgentsHandler
	req, _ = http.NewRequest("GET", "/api/options/agents", nil)
	rr = httptest.NewRecorder()
	http.HandlerFunc(GetAgentsHandler).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("GetAgentsHandler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
	var agents []string
	if err := json.NewDecoder(rr.Body).Decode(&agents); err != nil {
		t.Errorf("GetAgentsHandler returned invalid JSON: %v", err)
	}
	// Due to maps not having guaranteed order and tests running concurrently,
	// we just check that the list is populated.
	if len(agents) == 0 {
		t.Errorf("GetAgentsHandler returned empty list")
	}

	// Test GetSecurityAnalyzersHandler
	req, _ = http.NewRequest("GET", "/api/options/security-analyzers", nil)
	rr = httptest.NewRecorder()
	http.HandlerFunc(GetSecurityAnalyzersHandler).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("GetSecurityAnalyzersHandler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
	var analyzers []string
	if err := json.NewDecoder(rr.Body).Decode(&analyzers); err != nil {
		t.Errorf("GetSecurityAnalyzersHandler returned invalid JSON: %v", err)
	}
	if len(analyzers) == 0 {
		t.Errorf("GetSecurityAnalyzersHandler returned empty list")
	}
}
