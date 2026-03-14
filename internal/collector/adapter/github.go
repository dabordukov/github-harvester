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

	"google.golang.org/grpc/metadata"
)

const githubAPIEndpoint = "https://api.github.com"

var extractLastPageNumberRegEx = regexp.MustCompile(`page=(\d+)>; rel="last"`)

type RepositoryDTO struct {
	Name        string `json:"name"`
	OwnerStruct struct {
		Login string `json:"login"`
	} `json:"owner"`
	Description  string `json:"description"`
	Stars        int    `json:"stargazers_count"`
	Forks        int    `json:"forks"`
	CreatedAt    string `json:"created_at"`
	CommitsCount int
}

type GitHubAdapter struct {
	httpClient *http.Client
}

func NewGitHubAdapter() *GitHubAdapter {
	return &GitHubAdapter{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (ga *GitHubAdapter) FetchAll(ctx context.Context, owner, repoName string) (*RepositoryDTO, error) {
	repoStruct := &RepositoryDTO{}
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

func (ga *GitHubAdapter) GetRepoInfo(ctx context.Context, owner, repoName string, repoStruct *RepositoryDTO) error {
	url := fmt.Sprintf("%s/repos/%s/%s", githubAPIEndpoint, owner, repoName)
	req, err := http.NewRequest("GET", url, nil)
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

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
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

	err = json.NewDecoder(resp.Body).Decode(&repoStruct)
	if err != nil {
		return fmt.Errorf("can't parse JSON: %w", err)
	}

	return nil
}

func (ga *GitHubAdapter) GetCommitsCount(ctx context.Context, owner, repoName string, repoStruct *RepositoryDTO) error {
	url := fmt.Sprintf("%s/repos/%s/%s/commits?per_page=1", githubAPIEndpoint, owner, repoName)
	req, err := http.NewRequest("HEAD", url, nil)
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github api: %s", resp.Status)
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
