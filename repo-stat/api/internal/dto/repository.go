package dto

type RepositoryInfoResponse struct {
	Name         string `json:"name"`
	Owner        string `json:"owner"`
	Description  string `json:"description"`
	Forks        int64  `json:"forks"`
	Stars        int64  `json:"stars"`
	CreatedAt    string `json:"created_at"`
	CommitsCount int64  `json:"commits_count"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
