package usecase

import (
	"context"

	"repo-stat/api/internal/domain"
)

type RepositoryGetter interface {
	GetRepository(ctx context.Context, owner, repo string) (*domain.RepositoryInfo, error)
}

type RepositoryInfo struct {
	getter RepositoryGetter
}

func NewRepositoryInfo(getter RepositoryGetter) *RepositoryInfo {
	return &RepositoryInfo{getter: getter}
}

func (r *RepositoryInfo) Fetch(ctx context.Context, owner, repo string) (*domain.RepositoryInfo, error) {
	return r.getter.GetRepository(ctx, owner, repo)
}
