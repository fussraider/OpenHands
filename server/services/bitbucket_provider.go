package services

import (
	"context"
	"errors"
	"openhands-go/server/microagent"
)

type BitbucketProvider struct {
	BaseURL    string
	DataCenter bool
}

func NewBitbucketProvider(baseURL string, dataCenter bool) *BitbucketProvider {
	if baseURL == "" {
		if dataCenter {
			baseURL = "http://localhost:7990" // Default Bitbucket Data Center
		} else {
			baseURL = "https://api.bitbucket.org/2.0"
		}
	}
	return &BitbucketProvider{
		BaseURL:    baseURL,
		DataCenter: dataCenter,
	}
}

func (p *BitbucketProvider) ListRepositories(ctx context.Context, token string, page, perPage int, sort string) ([]GitRepository, error) {
	return nil, errors.New("BitbucketProvider.ListRepositories not implemented")
}

func (p *BitbucketProvider) SearchRepositories(ctx context.Context, token, query string, page, perPage int, sort, order string) ([]GitRepository, error) {
	return nil, errors.New("BitbucketProvider.SearchRepositories not implemented")
}

func (p *BitbucketProvider) GetBranches(ctx context.Context, token, owner, repo string, page, perPage int) ([]GitBranch, error) {
	return nil, errors.New("BitbucketProvider.GetBranches not implemented")
}

func (p *BitbucketProvider) SearchBranches(ctx context.Context, token, owner, repo, query string) ([]GitBranch, error) {
	return nil, errors.New("BitbucketProvider.SearchBranches not implemented")
}

func (p *BitbucketProvider) GetUser(ctx context.Context, token string) (*GitUser, error) {
	return nil, errors.New("BitbucketProvider.GetUser not implemented")
}

func (p *BitbucketProvider) GetInstallations(ctx context.Context, token string) ([]GitInstallation, error) {
	return nil, errors.New("BitbucketProvider.GetInstallations not implemented")
}

func (p *BitbucketProvider) GetFileContent(ctx context.Context, token, owner, repo, path string) (string, error) {
	return "", errors.New("BitbucketProvider.GetFileContent not implemented")
}

func (p *BitbucketProvider) GetMicroagents(ctx context.Context, token, owner, repo string) ([]microagent.MicroagentResponse, error) {
	return nil, errors.New("BitbucketProvider.GetMicroagents not implemented")
}

func (p *BitbucketProvider) GetMicroagentContent(ctx context.Context, token, owner, repo, path string) (*microagent.MicroagentContentResponse, error) {
	return nil, errors.New("BitbucketProvider.GetMicroagentContent not implemented")
}

func (p *BitbucketProvider) GetSuggestedTasks(ctx context.Context, token string) ([]SuggestedTask, error) {
	return nil, errors.New("BitbucketProvider.GetSuggestedTasks not implemented")
}

// CreatePR implements the MCP PR creation
func (p *BitbucketProvider) CreatePR(ctx context.Context, token string, repoName, sourceBranch, targetBranch, title, description string, draft bool, labels []string) (string, error) {
	return "", errors.New("BitbucketProvider.CreatePR not implemented")
}

func (p *BitbucketProvider) CreateMR(ctx context.Context, token string, id interface{}, sourceBranch, targetBranch, title, description string, labels []string) (string, error) {
	return "", errors.New("BitbucketProvider.CreateMR not implemented")
}
