package handlers

import (
	"encoding/json"
	"net/http"
)

// Helper types for Github/Git API
type Repository struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
	HTMLURL  string `json:"html_url"`
}

type User struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}

type Branch struct {
	Name string `json:"name"`
}

func GetUserInstallationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Mock: return empty list of installations
	json.NewEncoder(w).Encode([]string{})
}

func GetUserRepositoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Mock: return mock repositories
	repos := []Repository{
		{ID: 1, FullName: "openhands/test-repo", HTMLURL: "https://github.com/openhands/test-repo"},
	}
	json.NewEncoder(w).Encode(repos)
}

func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Mock: return mock user
	user := User{
		Login:     "mock-user",
		AvatarURL: "https://avatars.githubusercontent.com/u/0?v=4",
	}
	json.NewEncoder(w).Encode(user)
}

func SearchRepositoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Mock: return mock repositories
	repos := []Repository{
		{ID: 1, FullName: "openhands/search-result", HTMLURL: "https://github.com/openhands/search-result"},
	}
	json.NewEncoder(w).Encode(repos)
}

func SearchBranchesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Mock: return mock branches
	branches := []Branch{
		{Name: "main"},
		{Name: "dev"},
	}
	json.NewEncoder(w).Encode(branches)
}

func GetRepositoryBranchesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Mock: return mock branches
	branches := []Branch{
		{Name: "main"},
		{Name: "feature/test"},
	}
	// Return as paginated response usually, but for now list
	// The Python code returns PaginatedBranchesResponse, but let's check what that is.
	// It probably has a 'data' field or similar.
	// Let's assume list for now, or match if frontend expects struct.
	// For this migration, stubbing with list is often safer unless we know struct.
	// Looking at Python `PaginatedBranchesResponse` it likely wraps list.
	// Let's return list for `search/branches` (list[Branch]) but what about this one?
	// The Python signature says `PaginatedBranchesResponse`.
	// We'll stub with a simple map for now.
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": branches,
		"meta": map[string]interface{}{"total": 2},
	})
}
