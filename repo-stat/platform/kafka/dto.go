package kafka

type CollectRequest struct {
	RequestID string `json:"request_id"`
	Owner     string `json:"owner"`
	RepoName  string `json:"repo_name"`
}

type CollectResult struct {
	RequestID    string `json:"request_id"`
	Owner        string `json:"owner"`
	RepoName     string `json:"repo_name"`
	Description  string `json:"description"`
	Forks        int    `json:"forks"`
	Stars        int    `json:"stars"`
	CreatedAt    string `json:"created_at"`
	CommitsCount int    `json:"commits_count"`
	ErrorReason  string `json:"error_reason,omitempty"`
}
