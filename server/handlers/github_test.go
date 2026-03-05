package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"openhands-go/server/microagent"
	"openhands-go/server/services"
	"testing"
)

type MockGitProvider struct{}

func (m *MockGitProvider) ListRepositories(ctx context.Context, token string, page, perPage int, sort string) ([]services.GitRepository, error) {
	return []services.GitRepository{{ID: "1", FullName: "openhands/test-repo"}}, nil
}
func (m *MockGitProvider) SearchRepositories(ctx context.Context, token, query string, page, perPage int, sort, order string) ([]services.GitRepository, error) {
	return []services.GitRepository{{ID: "1", FullName: "openhands/search-result"}}, nil
}
func (m *MockGitProvider) GetBranches(ctx context.Context, token, owner, repo string, page, perPage int) ([]services.GitBranch, error) {
	return []services.GitBranch{{Name: "main"}}, nil
}
func (m *MockGitProvider) SearchBranches(ctx context.Context, token, owner, repo, query string) ([]services.GitBranch, error) {
	return []services.GitBranch{{Name: "main"}}, nil
}
func (m *MockGitProvider) GetUser(ctx context.Context, token string) (*services.GitUser, error) {
	return &services.GitUser{Login: "mock-user"}, nil
}
func (m *MockGitProvider) GetInstallations(ctx context.Context, token string) ([]services.GitInstallation, error) {
	return []services.GitInstallation{}, nil
}
func (m *MockGitProvider) GetFileContent(ctx context.Context, token, owner, repo, path string) (string, error) {
	return "mock content", nil
}
func (m *MockGitProvider) GetMicroagents(ctx context.Context, token, owner, repo string) ([]microagent.MicroagentResponse, error) {
	return []microagent.MicroagentResponse{}, nil
}
func (m *MockGitProvider) GetMicroagentContent(ctx context.Context, token, owner, repo, path string) (*microagent.MicroagentContentResponse, error) {
	return &microagent.MicroagentContentResponse{Content: "mock"}, nil
}
func (m *MockGitProvider) GetSuggestedTasks(ctx context.Context, token string) ([]services.SuggestedTask, error) {
	return []services.SuggestedTask{
		{
			GitProvider: "github",
			IssueNumber: 1,
			Repo:        "openhands/test-repo",
			Title:       "Test issue",
			TaskType:    "OPEN_ISSUE",
		},
	}, nil
}

func init() {
	GitService = &services.GitService{
		Provider: &MockGitProvider{},
	}
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
	var repos []services.GitRepository
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
	var user services.GitUser
	if err := json.NewDecoder(rr.Body).Decode(&user); err != nil {
		t.Errorf("GetUserInfoHandler returned invalid JSON: %v", err)
	}
	if user.Login == "" {
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

	// Test Microagent Handlers
	req, _ = http.NewRequest("GET", "/api/user/repository/owner/repo/microagents", nil)
	req.Header.Set("Authorization", tokenHeader)
	req.SetPathValue("owner", "owner")
	req.SetPathValue("repo", "repo")

	rr = httptest.NewRecorder()
	http.HandlerFunc(GetRepositoryMicroagentsHandler).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("GetRepositoryMicroagentsHandler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	req, _ = http.NewRequest("GET", "/api/user/repository/owner/repo/microagents/content?file_path=agent.md", nil)
	req.Header.Set("Authorization", tokenHeader)
	req.SetPathValue("owner", "owner")
	req.SetPathValue("repo", "repo")
	rr = httptest.NewRecorder()
	http.HandlerFunc(GetRepositoryMicroagentContentHandler).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("GetRepositoryMicroagentContentHandler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
}

func TestGetSuggestedTasksHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/user/suggested-tasks", nil)
	req.Header.Set("Authorization", "Bearer token")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetSuggestedTasksHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var tasks []services.SuggestedTask
	if err := json.NewDecoder(rr.Body).Decode(&tasks); err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 1 {
		t.Fatalf("expected 1 suggested task, got %d", len(tasks))
	}

	if tasks[0].Title != "Test issue" {
		t.Errorf("expected task title 'Test issue', got '%s'", tasks[0].Title)
	}
}
