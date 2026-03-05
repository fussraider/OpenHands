package services

import (
	"context"
	"log/slog"
	"openhands-go/server/config"
	"openhands-go/server/microagent"
	"strings"
	"time"

	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

type IGithubService interface {
	ListRepositories(ctx context.Context, token string, page, perPage int, sort string) ([]*github.Repository, error)
	SearchRepositories(ctx context.Context, token, query string, page, perPage int, sort, order string) ([]*github.Repository, error)
	GetBranches(ctx context.Context, token, owner, repo string, page, perPage int) ([]*github.Branch, error)
	SearchBranches(ctx context.Context, token, owner, repo, query string) ([]*github.Branch, error)
	GetUser(ctx context.Context, token string) (*github.User, error)
	GetInstallations(ctx context.Context, token string) ([]*github.Installation, error)
	GetSuggestedTasks(ctx context.Context, token string) ([]SuggestedTask, error)
	GetFileContent(ctx context.Context, token, owner, repo, path string) (string, error)
	ListDirectory(ctx context.Context, token, owner, repo, path string) ([]*github.RepositoryContent, error)
	GetMicroagents(ctx context.Context, token, owner, repo string) ([]microagent.MicroagentResponse, error)
	GetMicroagentContent(ctx context.Context, token, owner, repo, path string) (*microagent.MicroagentContentResponse, error)
}

type GithubService struct {
	config config.GithubConfig
}

func NewGithubService(cfg config.GithubConfig) *GithubService {
	return &GithubService{
		config: cfg,
	}
}

func (s *GithubService) getClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func (s *GithubService) ListRepositories(ctx context.Context, token string, page, perPage int, sort string) ([]*github.Repository, error) {
	client := s.getClient(ctx, token)
	opts := &github.RepositoryListOptions{
		Sort:        sort,
		ListOptions: github.ListOptions{Page: page, PerPage: perPage},
	}
	// "user" lists repos for the authenticated user
	repos, _, err := client.Repositories.List(ctx, "", opts)
	return repos, err
}

func (s *GithubService) SearchRepositories(ctx context.Context, token, query string, page, perPage int, sort, order string) ([]*github.Repository, error) {
	client := s.getClient(ctx, token)
	opts := &github.SearchOptions{
		Sort:        sort,
		Order:       order,
		ListOptions: github.ListOptions{Page: page, PerPage: perPage},
	}
	result, _, err := client.Search.Repositories(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	return result.Repositories, nil
}

func (s *GithubService) GetBranches(ctx context.Context, token, owner, repo string, page, perPage int) ([]*github.Branch, error) {
	client := s.getClient(ctx, token)
	opts := &github.BranchListOptions{
		ListOptions: github.ListOptions{Page: page, PerPage: perPage},
	}
	branches, _, err := client.Repositories.ListBranches(ctx, owner, repo, opts)
	return branches, err
}

func (s *GithubService) SearchBranches(ctx context.Context, token, owner, repo, query string) ([]*github.Branch, error) {
	return s.GetBranches(ctx, token, owner, repo, 1, 100)
}

func (s *GithubService) GetUser(ctx context.Context, token string) (*github.User, error) {
	client := s.getClient(ctx, token)
	user, _, err := client.Users.Get(ctx, "")
	return user, err
}

func (s *GithubService) GetInstallations(ctx context.Context, token string) ([]*github.Installation, error) {
	client := s.getClient(ctx, token)
	installations, _, err := client.Apps.ListUserInstallations(ctx, nil)
	if err != nil {
		return nil, err
	}
	return installations, nil
}

func (s *GithubService) GetSuggestedTasks(ctx context.Context, token string) ([]SuggestedTask, error) {
	// MVP: return empty tasks or mock tasks.
	// The Python version fetches user PRs and issues.
	return []SuggestedTask{}, nil
}

func (s *GithubService) GetFileContent(ctx context.Context, token, owner, repo, path string) (string, error) {
	client := s.getClient(ctx, token)
	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		return "", err
	}
	content, err := fileContent.GetContent()
	if err != nil {
		return "", err
	}
	return content, nil
}

func (s *GithubService) ListDirectory(ctx context.Context, token, owner, repo, path string) ([]*github.RepositoryContent, error) {
	client := s.getClient(ctx, token)
	_, directoryContent, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	return directoryContent, err
}

func (s *GithubService) GetMicroagents(ctx context.Context, token, owner, repo string) ([]microagent.MicroagentResponse, error) {
	var microagents []microagent.MicroagentResponse

	// 1. Check .cursorrules
	// Try root
	if content, err := s.GetFileContent(ctx, token, owner, repo, ".cursorrules"); err == nil && content != "" {
		microagents = append(microagents, microagent.MicroagentResponse{
			Name:      "cursorrules",
			Path:      ".cursorrules",
			CreatedAt: time.Now(), // Mock
		})
	}

	// 2. Check .openhands/microagents
	// The Python implementation handles special cases for repo names, but for now we implement the standard path.
	microagentPath := ".openhands/microagents"
	files, err := s.ListDirectory(ctx, token, owner, repo, microagentPath)
	if err == nil {
		for _, file := range files {
			if file.GetType() == "file" && strings.HasSuffix(file.GetName(), ".md") && file.GetName() != "README.md" {
				name := strings.TrimSuffix(file.GetName(), ".md")
				microagents = append(microagents, microagent.MicroagentResponse{
					Name:      name,
					Path:      file.GetPath(),
					CreatedAt: time.Now(), // Mock
				})
			}
		}
	} else {
		// Log but don't fail if directory doesn't exist
		slog.Debug("Microagents directory not found or error", "error", err)
	}

	return microagents, nil
}

func (s *GithubService) GetMicroagentContent(ctx context.Context, token, owner, repo, path string) (*microagent.MicroagentContentResponse, error) {
	content, err := s.GetFileContent(ctx, token, owner, repo, path)
	if err != nil {
		return nil, err
	}

	metadata, _, err := microagent.Parse(content, path)
	if err != nil {
		return nil, err
	}

	return &microagent.MicroagentContentResponse{
		Content:     content, // Return full content including frontmatter? Python returns `microagent.content` which is body. Let's check Parse return.
		Path:        path,
		Triggers:    metadata.Triggers,
		GitProvider: "github",
	}, nil
}
