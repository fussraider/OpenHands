package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/v60/github"
)

type MockGithubService struct{}

func (m *MockGithubService) ListRepositories(ctx context.Context, token string, page, perPage int, sort string) ([]*github.Repository, error) {
	return []*github.Repository{{ID: github.Int64(1), FullName: github.String("openhands/test-repo")}}, nil
}
func (m *MockGithubService) SearchRepositories(ctx context.Context, token, query string, page, perPage int, sort, order string) ([]*github.Repository, error) {
	return []*github.Repository{{ID: github.Int64(1), FullName: github.String("openhands/search-result")}}, nil
}
func (m *MockGithubService) GetBranches(ctx context.Context, token, owner, repo string, page, perPage int) ([]*github.Branch, error) {
	return []*github.Branch{{Name: github.String("main")}}, nil
}
func (m *MockGithubService) SearchBranches(ctx context.Context, token, owner, repo, query string) ([]*github.Branch, error) {
	return []*github.Branch{{Name: github.String("main")}}, nil
}
func (m *MockGithubService) GetUser(ctx context.Context, token string) (*github.User, error) {
	return &github.User{Login: github.String("mock-user")}, nil
}
func (m *MockGithubService) GetInstallations(ctx context.Context, token string) ([]*github.Installation, error) {
	return []*github.Installation{}, nil
}
func (m *MockGithubService) GetSuggestedTasks(ctx context.Context, token string) ([]interface{}, error) {
	return []interface{}{}, nil
}

func init() {
	GithubService = &MockGithubService{}
}

func TestGithubHandlers(t *testing.T) {
	// Need to set token header for handlers to work
	tokenHeader := "Bearer mock-token"

	// Test GetUserInstallationsHandler
	req, _ := http.NewRequest("GET", "/api/user/installations", nil)
	req.Header.Set("Authorization", tokenHeader)
	rr := httptest.NewRecorder()
	http.HandlerFunc(GetUserInstallationsHandler).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("GetUserInstallationsHandler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	// Test GetUserRepositoriesHandler
	req, _ = http.NewRequest("GET", "/api/user/repositories", nil)
	req.Header.Set("Authorization", tokenHeader)
	rr = httptest.NewRecorder()
	http.HandlerFunc(GetUserRepositoriesHandler).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("GetUserRepositoriesHandler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
	var repos []*github.Repository
	if err := json.NewDecoder(rr.Body).Decode(&repos); err != nil {
		t.Errorf("GetUserRepositoriesHandler returned invalid JSON: %v", err)
	}
	if len(repos) == 0 {
		t.Errorf("GetUserRepositoriesHandler returned empty list")
	}

	// Test GetUserInfoHandler
	req, _ = http.NewRequest("GET", "/api/user/info", nil)
	req.Header.Set("Authorization", tokenHeader)
	rr = httptest.NewRecorder()
	http.HandlerFunc(GetUserInfoHandler).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("GetUserInfoHandler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
	var user github.User
	if err := json.NewDecoder(rr.Body).Decode(&user); err != nil {
		t.Errorf("GetUserInfoHandler returned invalid JSON: %v", err)
	}
	if user.GetLogin() == "" {
		t.Errorf("GetUserInfoHandler returned empty login")
	}

	// Test SearchRepositoriesHandler
	req, _ = http.NewRequest("GET", "/api/user/search/repositories?query=test", nil)
	req.Header.Set("Authorization", tokenHeader)
	rr = httptest.NewRecorder()
	http.HandlerFunc(SearchRepositoriesHandler).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("SearchRepositoriesHandler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	// Test SearchBranchesHandler
	req, _ = http.NewRequest("GET", "/api/user/search/branches?repository=owner/repo&query=main", nil)
	req.Header.Set("Authorization", tokenHeader)
	rr = httptest.NewRecorder()
	http.HandlerFunc(SearchBranchesHandler).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("SearchBranchesHandler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
}
