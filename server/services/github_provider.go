package services

import (
	"context"
	"fmt"
	"log/slog"
	"openhands-go/server/config"
	"openhands-go/server/microagent"
	"strings"
	"time"

	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

type GitHubProvider struct {
	config config.GithubConfig
}

func NewGitHubProvider(cfg config.GithubConfig) *GitHubProvider {
	return &GitHubProvider{
		config: cfg,
	}
}

func (s *GitHubProvider) getClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func (s *GitHubProvider) ListRepositories(ctx context.Context, token string, page, perPage int, sort string) ([]GitRepository, error) {
	client := s.getClient(ctx, token)
	opts := &github.RepositoryListOptions{
		Sort:        sort,
		ListOptions: github.ListOptions{Page: page, PerPage: perPage},
	}
	repos, _, err := client.Repositories.List(ctx, "", opts)
	if err != nil {
		return nil, err
	}

	var result []GitRepository
	for _, r := range repos {
		result = append(result, GitRepository{
			ID:       fmt.Sprintf("%d", r.GetID()),
			FullName: r.GetFullName(),
			Private:  r.GetPrivate(),
		})
	}
	return result, nil
}

func (s *GitHubProvider) SearchRepositories(ctx context.Context, token, query string, page, perPage int, sort, order string) ([]GitRepository, error) {
	client := s.getClient(ctx, token)
	opts := &github.SearchOptions{
		Sort:        sort,
		Order:       order,
		ListOptions: github.ListOptions{Page: page, PerPage: perPage},
	}
	res, _, err := client.Search.Repositories(ctx, query, opts)
	if err != nil {
		return nil, err
	}

	var result []GitRepository
	for _, r := range res.Repositories {
		result = append(result, GitRepository{
			ID:       fmt.Sprintf("%d", r.GetID()),
			FullName: r.GetFullName(),
			Private:  r.GetPrivate(),
		})
	}
	return result, nil
}

func (s *GitHubProvider) GetBranches(ctx context.Context, token, owner, repo string, page, perPage int) ([]GitBranch, error) {
	client := s.getClient(ctx, token)
	opts := &github.BranchListOptions{
		ListOptions: github.ListOptions{Page: page, PerPage: perPage},
	}
	branches, _, err := client.Repositories.ListBranches(ctx, owner, repo, opts)
	if err != nil {
		return nil, err
	}

	var result []GitBranch
	for _, b := range branches {
		sha := ""
		if b.Commit != nil {
			sha = b.Commit.GetSHA()
		}
		result = append(result, GitBranch{
			Name: b.GetName(),
			SHA:  sha,
		})
	}
	return result, nil
}

func (s *GitHubProvider) SearchBranches(ctx context.Context, token, owner, repo, query string) ([]GitBranch, error) {
	// GitHub API doesn't have direct branch search, usually we filter list.
	// For simplicity, reuse GetBranches (ignoring query for now or filtering in memory)
	branches, err := s.GetBranches(ctx, token, owner, repo, 1, 100)
	if err != nil {
		return nil, err
	}
	// Simple filter
	if query == "" {
		return branches, nil
	}
	var filtered []GitBranch
	for _, b := range branches {
		if strings.Contains(b.Name, query) {
			filtered = append(filtered, b)
		}
	}
	return filtered, nil
}

func (s *GitHubProvider) GetUser(ctx context.Context, token string) (*GitUser, error) {
	client := s.getClient(ctx, token)
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}
	return &GitUser{
		Login:     user.GetLogin(),
		AvatarURL: user.GetAvatarURL(),
	}, nil
}

func (s *GitHubProvider) GetInstallations(ctx context.Context, token string) ([]GitInstallation, error) {
	client := s.getClient(ctx, token)
	installations, _, err := client.Apps.ListUserInstallations(ctx, nil)
	if err != nil {
		return nil, err
	}
	var result []GitInstallation
	for _, i := range installations {
		result = append(result, GitInstallation{
			ID: i.GetID(),
		})
	}
	return result, nil
}

func (s *GitHubProvider) GetFileContent(ctx context.Context, token, owner, repo, path string) (string, error) {
	client := s.getClient(ctx, token)
	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		return "", err
	}
	return fileContent.GetContent()
}

func (s *GitHubProvider) ListDirectory(ctx context.Context, token, owner, repo, path string) ([]*github.RepositoryContent, error) {
	client := s.getClient(ctx, token)
	_, directoryContent, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	return directoryContent, err
}

func (s *GitHubProvider) GetMicroagents(ctx context.Context, token, owner, repo string) ([]microagent.MicroagentResponse, error) {
	var microagents []microagent.MicroagentResponse

	// 1. Check .cursorrules
	if content, err := s.GetFileContent(ctx, token, owner, repo, ".cursorrules"); err == nil && content != "" {
		microagents = append(microagents, microagent.MicroagentResponse{
			Name:      "cursorrules",
			Path:      ".cursorrules",
			CreatedAt: time.Now(),
		})
	}

	// 2. Check .openhands/microagents
	microagentPath := ".openhands/microagents"
	files, err := s.ListDirectory(ctx, token, owner, repo, microagentPath)
	if err == nil {
		for _, file := range files {
			if file.GetType() == "file" && strings.HasSuffix(file.GetName(), ".md") && file.GetName() != "README.md" {
				name := strings.TrimSuffix(file.GetName(), ".md")
				microagents = append(microagents, microagent.MicroagentResponse{
					Name:      name,
					Path:      file.GetPath(),
					CreatedAt: time.Now(),
				})
			}
		}
	} else {
		slog.Debug("Microagents directory not found or error", "error", err)
	}

	return microagents, nil
}

func (s *GitHubProvider) GetMicroagentContent(ctx context.Context, token, owner, repo, path string) (*microagent.MicroagentContentResponse, error) {
	content, err := s.GetFileContent(ctx, token, owner, repo, path)
	if err != nil {
		return nil, err
	}

	metadata, _, err := microagent.Parse(content, path)
	if err != nil {
		return nil, err
	}

	return &microagent.MicroagentContentResponse{
		Content:     content,
		Path:        path,
		Triggers:    metadata.Triggers,
		GitProvider: "github",
	}, nil
}

func (s *GitHubProvider) GetSuggestedTasks(ctx context.Context, token string) ([]SuggestedTask, error) {
	client := s.getClient(ctx, token)
	opts := &github.IssueListOptions{
		State:       "open",
		Filter:      "assigned",
		ListOptions: github.ListOptions{PerPage: 15},
	}
	issues, _, err := client.Issues.List(ctx, true, opts)
	if err != nil {
		return nil, err
	}
	var tasks []SuggestedTask
	for _, issue := range issues {
		repoURL := issue.GetRepositoryURL()
		// repoURL looks like https://api.github.com/repos/owner/repo
		parts := strings.Split(repoURL, "/")
		repo := ""
		if len(parts) >= 2 {
			repo = parts[len(parts)-2] + "/" + parts[len(parts)-1]
		}
		taskType := "OPEN_ISSUE"
		if issue.IsPullRequest() {
			taskType = "PULL_REQUEST"
		}
		tasks = append(tasks, SuggestedTask{
			GitProvider: "github",
			IssueNumber: issue.GetNumber(),
			Repo:        repo,
			Title:       issue.GetTitle(),
			TaskType:    taskType,
		})
	}
	return tasks, nil
}
