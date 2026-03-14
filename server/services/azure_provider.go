package services

import (
	"context"
	"errors"
	"openhands-go/server/microagent"
)

type AzureDevOpsProvider struct {
	BaseURL string
}

func NewAzureDevOpsProvider(baseURL string) *AzureDevOpsProvider {
	if baseURL == "" {
		baseURL = "https://dev.azure.com"
	}
	return &AzureDevOpsProvider{
		BaseURL: baseURL,
	}
}

func (p *AzureDevOpsProvider) ListRepositories(ctx context.Context, token string, page, perPage int, sort string) ([]GitRepository, error) {
	return nil, errors.New("AzureDevOpsProvider.ListRepositories not implemented")
}

func (p *AzureDevOpsProvider) SearchRepositories(ctx context.Context, token, query string, page, perPage int, sort, order string) ([]GitRepository, error) {
	return nil, errors.New("AzureDevOpsProvider.SearchRepositories not implemented")
}

func (p *AzureDevOpsProvider) GetBranches(ctx context.Context, token, owner, repo string, page, perPage int) ([]GitBranch, error) {
	return nil, errors.New("AzureDevOpsProvider.GetBranches not implemented")
}

func (p *AzureDevOpsProvider) SearchBranches(ctx context.Context, token, owner, repo, query string) ([]GitBranch, error) {
	return nil, errors.New("AzureDevOpsProvider.SearchBranches not implemented")
}

func (p *AzureDevOpsProvider) GetUser(ctx context.Context, token string) (*GitUser, error) {
	return nil, errors.New("AzureDevOpsProvider.GetUser not implemented")
}

func (p *AzureDevOpsProvider) GetInstallations(ctx context.Context, token string) ([]GitInstallation, error) {
	return nil, errors.New("AzureDevOpsProvider.GetInstallations not implemented")
}

func (p *AzureDevOpsProvider) GetFileContent(ctx context.Context, token, owner, repo, path string) (string, error) {
	return "", errors.New("AzureDevOpsProvider.GetFileContent not implemented")
}

func (p *AzureDevOpsProvider) GetMicroagents(ctx context.Context, token, owner, repo string) ([]microagent.MicroagentResponse, error) {
	return nil, errors.New("AzureDevOpsProvider.GetMicroagents not implemented")
}

func (p *AzureDevOpsProvider) GetMicroagentContent(ctx context.Context, token, owner, repo, path string) (*microagent.MicroagentContentResponse, error) {
	return nil, errors.New("AzureDevOpsProvider.GetMicroagentContent not implemented")
}

func (p *AzureDevOpsProvider) GetSuggestedTasks(ctx context.Context, token string) ([]SuggestedTask, error) {
	return nil, errors.New("AzureDevOpsProvider.GetSuggestedTasks not implemented")
}

// CreatePR implements the MCP PR creation
func (p *AzureDevOpsProvider) CreatePR(ctx context.Context, token string, repoName, sourceBranch, targetBranch, title, description string, draft bool, labels []string) (string, error) {
	return "", errors.New("AzureDevOpsProvider.CreatePR not implemented")
}

func (p *AzureDevOpsProvider) CreateMR(ctx context.Context, token string, id interface{}, sourceBranch, targetBranch, title, description string, labels []string) (string, error) {
	return "", errors.New("AzureDevOpsProvider.CreateMR not implemented")
}
