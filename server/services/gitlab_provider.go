package services

import (
	"context"
	"errors"
	"openhands-go/server/microagent"
)

type GitLabProvider struct {
	BaseURL string
}

func NewGitLabProvider(baseURL string) *GitLabProvider {
	if baseURL == "" {
		baseURL = "https://gitlab.com"
	}
	return &GitLabProvider{
		BaseURL: baseURL,
	}
}

func (p *GitLabProvider) ListRepositories(ctx context.Context, token string, page, perPage int, sort string) ([]GitRepository, error) {
	return nil, errors.New("GitLabProvider.ListRepositories not implemented")
}

func (p *GitLabProvider) SearchRepositories(ctx context.Context, token, query string, page, perPage int, sort, order string) ([]GitRepository, error) {
	return nil, errors.New("GitLabProvider.SearchRepositories not implemented")
}

func (p *GitLabProvider) GetBranches(ctx context.Context, token, owner, repo string, page, perPage int) ([]GitBranch, error) {
	return nil, errors.New("GitLabProvider.GetBranches not implemented")
}

func (p *GitLabProvider) SearchBranches(ctx context.Context, token, owner, repo, query string) ([]GitBranch, error) {
	return nil, errors.New("GitLabProvider.SearchBranches not implemented")
}

func (p *GitLabProvider) GetUser(ctx context.Context, token string) (*GitUser, error) {
	return nil, errors.New("GitLabProvider.GetUser not implemented")
}

func (p *GitLabProvider) GetInstallations(ctx context.Context, token string) ([]GitInstallation, error) {
	return nil, errors.New("GitLabProvider.GetInstallations not implemented")
}

func (p *GitLabProvider) GetFileContent(ctx context.Context, token, owner, repo, path string) (string, error) {
	return "", errors.New("GitLabProvider.GetFileContent not implemented")
}

func (p *GitLabProvider) GetMicroagents(ctx context.Context, token, owner, repo string) ([]microagent.MicroagentResponse, error) {
	return nil, errors.New("GitLabProvider.GetMicroagents not implemented")
}

func (p *GitLabProvider) GetMicroagentContent(ctx context.Context, token, owner, repo, path string) (*microagent.MicroagentContentResponse, error) {
	return nil, errors.New("GitLabProvider.GetMicroagentContent not implemented")
}

func (p *GitLabProvider) GetSuggestedTasks(ctx context.Context, token string) ([]SuggestedTask, error) {
	return nil, errors.New("GitLabProvider.GetSuggestedTasks not implemented")
}

// CreateMR implements the MCP MR creation
func (p *GitLabProvider) CreateMR(ctx context.Context, token string, id interface{}, sourceBranch, targetBranch, title, description string, labels []string) (string, error) {
	return "", errors.New("GitLabProvider.CreateMR not implemented")
}

func (p *GitLabProvider) CreatePR(ctx context.Context, token string, repoName, sourceBranch, targetBranch, title, description string, draft bool, labels []string) (string, error) {
	return "", errors.New("GitLabProvider.CreatePR not implemented")
}
