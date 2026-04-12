package dto

type RepositoryDTO struct {
	FullName    string `json:"full_name"`
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
