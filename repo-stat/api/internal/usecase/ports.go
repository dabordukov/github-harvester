package usecase

import (
	"context"

	"repo-stat/api/internal/domain"
)

type Pinger interface {
	Ping(ctx context.Context) domain.PingStatus
}

type RepositoryGetter interface {
	GetRepository(ctx context.Context, owner, repo string) (*domain.RepositoryInfo, error)
	GetSubscriptionsInfo(ctx context.Context) ([]domain.RepositoryInfo, error)
}

type SubscriptionManager interface {
	CreateSubscription(ctx context.Context, owner, repo string) (*domain.Subscription, error)
	DeleteSubscription(ctx context.Context, owner, repo string) error
	ListSubscriptions(ctx context.Context) ([]domain.Subscription, error)
}
