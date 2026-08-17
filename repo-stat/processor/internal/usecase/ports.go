package usecase

import (
	"context"

	"repo-stat/processor/internal/domain"
)

type RepositoryStore interface {
	Get(ctx context.Context, owner, repo string) (*domain.Repository, error)
	Save(ctx context.Context, repository *domain.Repository) error
}

type KafkaPublisher interface {
	PublishCollectRequest(ctx context.Context, owner, repo string) error
}

type CollectorProvider interface {
	GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
	GetSubscriptionsInfo(ctx context.Context) ([]domain.Repository, error)
	Ping(ctx context.Context) string
}
