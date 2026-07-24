package adapter

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc/metadata"
)

const githubAPIEndpoint = "https://api.github.com"

type GitHubAdapter struct {
	httpClient *http.Client
}

func NewGitHubAdapter() *GitHubAdapter {
	return &GitHubAdapter{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (ga *GitHubAdapter) EnsureRepositoryExists(ctx context.Context, owner, repoName string) error {
	url := fmt.Sprintf("%s/repos/%s/%s", githubAPIEndpoint, owner, repoName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	if token := ga.extractToken(ctx); token != "" {
		req.Header.Set("Authorization", token)
	}

	resp, err := ga.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("can't close response body: %v\n", err)
		}
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrRateLimited
	default:
		return fmt.Errorf("github api unexpected status: %d", resp.StatusCode)
	}
}

func (ga *GitHubAdapter) extractToken(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if tokens := md.Get("authorization"); len(tokens) > 0 && tokens[0] != "" {
			return tokens[0]
		}
	}

	return ""
}
