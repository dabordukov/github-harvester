package model

type RepositoryModel struct {
	Name         string
	Owner        string
	Description  string
	Forks        int64
	Stars        int64
	CreatedAt    string
	CommitsCount int64
}
