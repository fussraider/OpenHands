package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGithubHandlers(t *testing.T) {
	// Test GetUserInstallationsHandler
	req, _ := http.NewRequest("GET", "/api/user/installations", nil)
	rr := httptest.NewRecorder()
	http.HandlerFunc(GetUserInstallationsHandler).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("GetUserInstallationsHandler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	// Test GetUserRepositoriesHandler
	req, _ = http.NewRequest("GET", "/api/user/repositories", nil)
	rr = httptest.NewRecorder()
	http.HandlerFunc(GetUserRepositoriesHandler).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("GetUserRepositoriesHandler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
	var repos []Repository
	if err := json.NewDecoder(rr.Body).Decode(&repos); err != nil {
		t.Errorf("GetUserRepositoriesHandler returned invalid JSON: %v", err)
	}
	if len(repos) == 0 {
		t.Errorf("GetUserRepositoriesHandler returned empty list")
	}

	// Test GetUserInfoHandler
	req, _ = http.NewRequest("GET", "/api/user/info", nil)
	rr = httptest.NewRecorder()
	http.HandlerFunc(GetUserInfoHandler).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("GetUserInfoHandler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
	var user User
	if err := json.NewDecoder(rr.Body).Decode(&user); err != nil {
		t.Errorf("GetUserInfoHandler returned invalid JSON: %v", err)
	}
	if user.Login == "" {
		t.Errorf("GetUserInfoHandler returned empty login")
	}

	// Test SearchRepositoriesHandler
	req, _ = http.NewRequest("GET", "/api/user/search/repositories", nil)
	rr = httptest.NewRecorder()
	http.HandlerFunc(SearchRepositoriesHandler).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("SearchRepositoriesHandler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	// Test SearchBranchesHandler
	req, _ = http.NewRequest("GET", "/api/user/search/branches", nil)
	rr = httptest.NewRecorder()
	http.HandlerFunc(SearchBranchesHandler).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("SearchBranchesHandler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
}
