package services

import (
	"context"
	"openhands-go/server/microagent"
)

type GitUser struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}

type GitRepository struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
}

type GitBranch struct {
	Name string `json:"name"`
	SHA  string `json:"sha"`
}

type GitInstallation struct {
	ID int64 `json:"id"`
}

type SuggestedTask struct {
	GitProvider string `json:"git_provider"`
	IssueNumber int    `json:"issue_number"`
	Repo        string `json:"repo"`
	Title       string `json:"title"`
	TaskType    string `json:"task_type"` // e.g. "OPEN_ISSUE", "MERGE_CONFLICTS", "FAILING_CHECKS", "UNRESOLVED_COMMENTS"
}

// GitProvider defines the common interface for git providers (GitHub, GitLab, etc.)
type GitProvider interface {
	ListRepositories(ctx context.Context, token string, page, perPage int, sort string) ([]GitRepository, error)
	SearchRepositories(ctx context.Context, token, query string, page, perPage int, sort, order string) ([]GitRepository, error)
	GetBranches(ctx context.Context, token, owner, repo string, page, perPage int) ([]GitBranch, error)
	SearchBranches(ctx context.Context, token, owner, repo, query string) ([]GitBranch, error)
	GetUser(ctx context.Context, token string) (*GitUser, error)
	GetInstallations(ctx context.Context, token string) ([]GitInstallation, error)
	GetFileContent(ctx context.Context, token, owner, repo, path string) (string, error)
	GetMicroagents(ctx context.Context, token, owner, repo string) ([]microagent.MicroagentResponse, error)
	GetMicroagentContent(ctx context.Context, token, owner, repo, path string) (*microagent.MicroagentContentResponse, error)
	GetSuggestedTasks(ctx context.Context, token string) ([]SuggestedTask, error)
	CreatePR(ctx context.Context, token string, repoName, sourceBranch, targetBranch, title, description string, draft bool, labels []string) (string, error)
	CreateMR(ctx context.Context, token string, id interface{}, sourceBranch, targetBranch, title, description string, labels []string) (string, error)
}
