package service

import (
	"context"

	"github-harvester/internal/collector/adapter"
)

type GitHubProvider interface {
	FetchAll(ctx context.Context, owner, repo string) (*adapter.RepositoryDTO, error)
}

type CollectorService struct {
	provider GitHubProvider
}

type RepositoryModel struct {
	Name         string
	Owner        string
	Description  string
	Forks        int64
	Stars        int64
	CreatedAt    string
	CommitsCount int64
}

func NewCollectorService(p GitHubProvider) *CollectorService {
	return &CollectorService{provider: p}
}

func (s *CollectorService) GetRepositoryData(ctx context.Context, owner, repo string) (*RepositoryModel, error) {
	data, err := s.provider.FetchAll(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	return &RepositoryModel{
		Name:         data.Name,
		Owner:        data.OwnerStruct.Login,
		Description:  data.Description,
		Forks:        int64(data.Forks),
		Stars:        int64(data.Stars),
		CreatedAt:    data.CreatedAt,
		CommitsCount: int64(data.CommitsCount),
	}, nil
}
