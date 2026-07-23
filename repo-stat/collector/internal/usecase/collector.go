package service

import (
	"context"
	"log/slog"

	"repo-stat/collector/internal/dto"
	subscriberpb "repo-stat/proto/subscriber"
)

type GitHubProvider interface {
	FetchAll(ctx context.Context, owner, repo string) (*dto.RepositoryDTO, error)
}

type SubscriptionProvider interface {
	ListSubscriptions(ctx context.Context) ([]*subscriberpb.Subscription, error)
}

type CollectorService struct {
	provider   GitHubProvider
	subscriber SubscriptionProvider
}

type RepositoryModel struct {
	FullName     string
	Name         string
	Owner        string
	Description  string
	Forks        int64
	Stars        int64
	CreatedAt    string
	CommitsCount int64
}

func NewCollectorService(p GitHubProvider, subscriber SubscriptionProvider) *CollectorService {
	return &CollectorService{
		provider:   p,
		subscriber: subscriber,
	}
}

func (s *CollectorService) GetRepositoryData(ctx context.Context, owner, repo string) (*RepositoryModel, error) {
	data, err := s.provider.FetchAll(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	return &RepositoryModel{
		FullName:     data.FullName,
		Name:         data.Name,
		Owner:        data.OwnerStruct.Login,
		Description:  data.Description,
		Forks:        int64(data.Forks),
		Stars:        int64(data.Stars),
		CreatedAt:    data.CreatedAt,
		CommitsCount: int64(data.CommitsCount),
	}, nil
}

func (s *CollectorService) GetSubscriptionsData(ctx context.Context) ([]RepositoryModel, error) {
	subscriptions, err := s.subscriber.ListSubscriptions(ctx)
	if err != nil {
		return nil, err
	}

	repositories := make([]RepositoryModel, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		repository, err := s.GetRepositoryData(ctx, subscription.GetOwner(), subscription.GetRepoName())
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		repositories = append(repositories, *repository)
	}

	return repositories, nil
}
