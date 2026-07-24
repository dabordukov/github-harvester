package domain

type RepositoryInfo struct {
	FullName     string
	Name         string
	Owner        string
	Description  string
	Forks        int64
	Stars        int64
	CreatedAt    string
	CommitsCount int64
}

type Subscription struct {
	Owner    string
	RepoName string
}
