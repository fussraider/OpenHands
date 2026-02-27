package services

import (
	"context"
	"openhands-go/server/config"

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
	GetSuggestedTasks(ctx context.Context, token string) ([]interface{}, error)
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

func (s *GithubService) GetSuggestedTasks(ctx context.Context, token string) ([]interface{}, error) {
	return []interface{}{}, nil
}
