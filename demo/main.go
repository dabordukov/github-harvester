package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
	"text/tabwriter"
	"time"
)

var statusMessages = map[int]string{
	http.StatusNotFound:            "Репозиторий не найден. Проверьте правильность owner/repo.",
	http.StatusUnauthorized:        "Ошибка авторизации. Проверьте GITHUB_TOKEN.",
	http.StatusForbidden:           "Доступ запрещен. Возможно, исчерпан лимит запросов (Rate Limit).",
	http.StatusInternalServerError: "Ошибка на стороне GitHub. Попробуйте позже.",
}

var token string = os.Getenv("GITHUB_TOKEN")

type GithubHarvester struct {
	apiEndpoint string
	repoStruct  *RepositoryInfo
	httpClient  *http.Client
	owner       string
	repoName    string
}

type Owner struct {
	Login string `json:"login"`
}

type RepositoryInfo struct {
	Name         string `json:"name"`
	Owner        Owner  `json:"owner"`
	Description  string `json:"description"`
	Forks        int    `json:"forks"`
	Stargazers   int    `json:"stargazers_count"`
	CreatedAt    string `json:"created_at"`
	CommitsCount int
}

func main() {
	var owner string
	var repoName string

	if len(os.Args) == 3 {
		owner, repoName = os.Args[1], os.Args[2]
	} else {
		fmt.Println("Enter github repository owner and repository name (separate with space)")
		if _, err := fmt.Scan(&owner, &repoName); err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			return
		}
	}

	harvester := NewGithubHarvester(owner, repoName)
	if err := harvester.HarvestAll(); err != nil {
		fmt.Println(err)
		return
	}

	harvester.PrintRepoInformation()
}

func NewGithubHarvester(owner, repoName string) *GithubHarvester {
	return &GithubHarvester{
		apiEndpoint: "https://api.github.com",
		repoStruct:  &RepositoryInfo{},
		httpClient:  &http.Client{Timeout: 10 * time.Second},
		owner:       owner,
		repoName:    repoName,
	}
}

func (gh *GithubHarvester) HarvestAll() error {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := gh.GetRepoInfo(); err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := gh.GetCommitsCount(); err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return <-errChan
	}

	return nil
}

func (gh *GithubHarvester) PrintRepoInformation() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	_, _ = fmt.Fprintln(w, "--------------------------------------------------")
	_, _ = fmt.Fprintf(w, "🚀 REPOSITORY:\t%s/%s\n", gh.repoStruct.Owner.Login, gh.repoStruct.Name)
	_, _ = fmt.Fprintln(w, "--------------------------------------------------")

	if gh.repoStruct.Description != "" {
		_, _ = fmt.Fprintf(w, "📝 Description:\t%s\n", gh.repoStruct.Description)
	} else {
		_, _ = fmt.Fprintf(w, "📝 Description:\t[No description provided]\n")
	}

	_, _ = fmt.Fprintf(w, "⭐ Stars:\t%d\n", gh.repoStruct.Stargazers)
	_, _ = fmt.Fprintf(w, "🍴 Forks:\t%d\n", gh.repoStruct.Forks)
	_, _ = fmt.Fprintf(w, "📦 Commits:\t%d\n", gh.repoStruct.CommitsCount)
	_, _ = fmt.Fprintf(w, "📅 Created:\t%s\n", gh.repoStruct.CreatedAt)
	_, _ = fmt.Fprintln(w, "--------------------------------------------------")

	if err := w.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "Can't flush buffer: %v", err)
	}
}

func (gh *GithubHarvester) GetRepoInfo() error {
	url := fmt.Sprintf("%s/repos/%s/%s", gh.apiEndpoint, gh.owner, gh.repoName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := gh.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Can't close response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return mapError(resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(gh.repoStruct)
	if err != nil {
		return fmt.Errorf("can't parse JSON: %w", err)
	}

	return nil
}

func (gh *GithubHarvester) GetCommitsCount() error {
	url := fmt.Sprintf("%s/repos/%s/%s/commits?per_page=1", gh.apiEndpoint, gh.owner, gh.repoName)
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return err
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := gh.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Can't close response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return mapError(resp.StatusCode)
	}

	link := resp.Header.Get("link")
	if link == "" {
		return fmt.Errorf("can't get commits")
	}

	re := regexp.MustCompile(`page=(\d+)>; rel="last"`)
	matches := re.FindStringSubmatch(link)

	if len(matches) > 1 {
		count, _ := strconv.Atoi(matches[1])
		gh.repoStruct.CommitsCount = count
	}

	return nil
}

func mapError(statusCode int) error {
	if msg, ok := statusMessages[statusCode]; ok {
		return fmt.Errorf("%d: %s", statusCode, msg)
	}

	return fmt.Errorf("ошибка %d", statusCode)
}
