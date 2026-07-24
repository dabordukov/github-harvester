package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	"repo-stat/collector/internal/dto"
	service "repo-stat/collector/internal/usecase"

	"google.golang.org/grpc/metadata"
)

const githubAPIEndpoint = "https://api.github.com"

var extractLastPageNumberRegEx = regexp.MustCompile(`page=(\d+)>; rel="last"`)

type GitHubAdapter struct {
	httpClient *http.Client
}

func NewGitHubAdapter() *GitHubAdapter {
	return &GitHubAdapter{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (ga *GitHubAdapter) FetchAll(ctx context.Context, owner, repoName string) (*dto.RepositoryDTO, error) {
	repoStruct := &dto.RepositoryDTO{}
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := ga.GetRepoInfo(ctx, owner, repoName, repoStruct); err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := ga.GetCommitsCount(ctx, owner, repoName, repoStruct); err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	return repoStruct, nil
}

func (ga *GitHubAdapter) GetRepoInfo(ctx context.Context, owner, repoName string, repoStruct *dto.RepositoryDTO) error {
	url := fmt.Sprintf("%s/repos/%s/%s", githubAPIEndpoint, owner, repoName)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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
			log.Printf("Can't close response body: %v\n", err)
		}
	}()

	if err := matchStatusCode(resp.StatusCode); err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&repoStruct)
	if err != nil {
		return fmt.Errorf("can't parse JSON: %w", err)
	}

	return nil
}

func (ga *GitHubAdapter) GetCommitsCount(ctx context.Context, owner, repoName string, repoStruct *dto.RepositoryDTO) error {
	url := fmt.Sprintf("%s/repos/%s/%s/commits?per_page=1", githubAPIEndpoint, owner, repoName)
	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
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
			log.Printf("Can't close response body: %v\n", err)
		}
	}()

	if err := matchStatusCode(resp.StatusCode); err != nil {
		return err
	}

	link := resp.Header.Get("link")
	if link == "" {
		return fmt.Errorf("can't get commits")
	}

	matches := extractLastPageNumberRegEx.FindStringSubmatch(link)

	if len(matches) > 1 {
		count, _ := strconv.Atoi(matches[1])
		repoStruct.CommitsCount = count
	}

	return nil
}

func (ga *GitHubAdapter) extractToken(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if tokens := md.Get("authorization"); len(tokens) > 0 {
			return tokens[0]
		}
	}
	return ""
}

func matchStatusCode(code int) error {
	switch code {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return service.ErrNotFound
	case http.StatusUnauthorized:
		return service.ErrUnauthorized
	case http.StatusForbidden:
		return service.ErrRateLimited
	default:
		return fmt.Errorf("%w: %d", service.ErrUnexpected, code)
	}
}
