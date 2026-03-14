package services

import (
	"context"
	"openhands-go/server/config"
	"openhands-go/server/microagent"
)

type GitService struct {
	Provider GitProvider // Exported for mocking in tests
}

func NewGitService(cfg *config.Config) *GitService {
	return &GitService{
		Provider: NewGitHubProvider(cfg.Github),
	}
}

func (s *GitService) GetProvider(ctx context.Context, token string) (GitProvider, error) {
	// Simple logic for determining provider based on context or token prefix
	// In a real implementation this should be passed in via API context
	// For MVP, if there is a way to determine it, we'll return the specific provider

	// Default to GitHub
	return s.Provider, nil
}

// Delegate methods to the appropriate provider
func (s *GitService) ListRepositories(ctx context.Context, token string, page, perPage int, sort string) ([]GitRepository, error) {
	p, _ := s.GetProvider(ctx, token)
	return p.ListRepositories(ctx, token, page, perPage, sort)
}

func (s *GitService) SearchRepositories(ctx context.Context, token, query string, page, perPage int, sort, order string) ([]GitRepository, error) {
	p, _ := s.GetProvider(ctx, token)
	return p.SearchRepositories(ctx, token, query, page, perPage, sort, order)
}

func (s *GitService) GetBranches(ctx context.Context, token, owner, repo string, page, perPage int) ([]GitBranch, error) {
	p, _ := s.GetProvider(ctx, token)
	return p.GetBranches(ctx, token, owner, repo, page, perPage)
}

func (s *GitService) SearchBranches(ctx context.Context, token, owner, repo, query string) ([]GitBranch, error) {
	p, _ := s.GetProvider(ctx, token)
	return p.SearchBranches(ctx, token, owner, repo, query)
}

func (s *GitService) GetUser(ctx context.Context, token string) (*GitUser, error) {
	p, _ := s.GetProvider(ctx, token)
	return p.GetUser(ctx, token)
}

func (s *GitService) GetInstallations(ctx context.Context, token string) ([]GitInstallation, error) {
	p, _ := s.GetProvider(ctx, token)
	return p.GetInstallations(ctx, token)
}

func (s *GitService) GetFileContent(ctx context.Context, token, owner, repo, path string) (string, error) {
	p, _ := s.GetProvider(ctx, token)
	return p.GetFileContent(ctx, token, owner, repo, path)
}

func (s *GitService) GetMicroagents(ctx context.Context, token, owner, repo string) ([]microagent.MicroagentResponse, error) {
	p, _ := s.GetProvider(ctx, token)
	return p.GetMicroagents(ctx, token, owner, repo)
}

func (s *GitService) GetMicroagentContent(ctx context.Context, token, owner, repo, path string) (*microagent.MicroagentContentResponse, error) {
	p, _ := s.GetProvider(ctx, token)
	return p.GetMicroagentContent(ctx, token, owner, repo, path)
}

// IGithubService compatibility (deprecated)
func (s *GitService) GetSuggestedTasks(ctx context.Context, token string) ([]SuggestedTask, error) {
	p, _ := s.GetProvider(ctx, token)
	if p == nil {
		return []SuggestedTask{}, nil
	}
	return p.GetSuggestedTasks(ctx, token)
}

func (s *GitService) CreatePR(ctx context.Context, token string, repoName, sourceBranch, targetBranch, title, description string, draft bool, labels []string) (string, error) {
	p, _ := s.GetProvider(ctx, token)
	return p.CreatePR(ctx, token, repoName, sourceBranch, targetBranch, title, description, draft, labels)
}

func (s *GitService) CreateMR(ctx context.Context, token string, id interface{}, sourceBranch, targetBranch, title, description string, labels []string) (string, error) {
	p, _ := s.GetProvider(ctx, token)
	return p.CreateMR(ctx, token, id, sourceBranch, targetBranch, title, description, labels)
}
