package service

import (
	"context"

	"repo-stat/processor/internal/domain"
)

type CollectorProvider interface {
	GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
	GetSubscriptionsInfo(ctx context.Context) ([]domain.Repository, error)
	Ping(ctx context.Context) string
}
