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
	// Simple logic: if token starts with "ghp_" or "github_pat_" or just default to GitHub for now.
	// Real implementation might need more explicit provider selection from request.
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
func (s *GitService) GetSuggestedTasks(ctx context.Context, token string) ([]interface{}, error) {
	return []interface{}{}, nil
}
