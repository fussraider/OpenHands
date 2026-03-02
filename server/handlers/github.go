package handlers

import (
	"encoding/json"
	"net/http"
	"openhands-go/server/services"
	"strconv"
	"strings"
)

var GitService *services.GitService

func getToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return r.Header.Get("X-Github-Token")
}

func GetUserInstallationsHandler(w http.ResponseWriter, r *http.Request) {
	token := getToken(r)
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	installations, err := GitService.GetInstallations(r.Context(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(installations)
}

func GetUserRepositoriesHandler(w http.ResponseWriter, r *http.Request) {
	token := getToken(r)
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 {
		perPage = 30
	}
	sort := r.URL.Query().Get("sort")

	repos, err := GitService.ListRepositories(r.Context(), token, page, perPage, sort)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repos)
}

func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	token := getToken(r)
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := GitService.GetUser(r.Context(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func SearchRepositoriesHandler(w http.ResponseWriter, r *http.Request) {
	token := getToken(r)
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := r.URL.Query().Get("query")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	repos, err := GitService.SearchRepositories(r.Context(), token, query, page, perPage, sort, order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repos)
}

func SearchBranchesHandler(w http.ResponseWriter, r *http.Request) {
	token := getToken(r)
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	repository := r.URL.Query().Get("repository") // owner/repo
	query := r.URL.Query().Get("query")

	parts := strings.Split(repository, "/")
	if len(parts) != 2 {
		http.Error(w, "Invalid repository format", http.StatusBadRequest)
		return
	}
	owner, repo := parts[0], parts[1]

	branches, err := GitService.SearchBranches(r.Context(), token, owner, repo, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(branches)
}

func GetRepositoryBranchesHandler(w http.ResponseWriter, r *http.Request) {
	token := getToken(r)
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	repository := r.URL.Query().Get("repository")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 30
	}

	parts := strings.Split(repository, "/")
	if len(parts) != 2 {
		http.Error(w, "Invalid repository format", http.StatusBadRequest)
		return
	}
	owner, repo := parts[0], parts[1]

	branches, err := GitService.GetBranches(r.Context(), token, owner, repo, page, perPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Python expects PaginatedBranchesResponse
	response := map[string]interface{}{
		"data": branches,
		// "meta": ...
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetRepositoryMicroagentsHandler(w http.ResponseWriter, r *http.Request) {
	token := getToken(r)
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	owner := r.PathValue("owner")
	repo := r.PathValue("repo")
	if owner == "" || repo == "" {
		http.Error(w, "Invalid repository format", http.StatusBadRequest)
		return
	}

	microagents, err := GitService.GetMicroagents(r.Context(), token, owner, repo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(microagents)
}

func GetRepositoryMicroagentContentHandler(w http.ResponseWriter, r *http.Request) {
	token := getToken(r)
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	owner := r.PathValue("owner")
	repo := r.PathValue("repo")
	filePath := r.URL.Query().Get("file_path")

	if filePath == "" {
		http.Error(w, "file_path parameter is required", http.StatusBadRequest)
		return
	}

	if owner == "" || repo == "" {
		http.Error(w, "Invalid repository format", http.StatusBadRequest)
		return
	}

	content, err := GitService.GetMicroagentContent(r.Context(), token, owner, repo, filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}
