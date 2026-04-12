package usecase

import (
	"context"

	"repo-stat/api/internal/domain"
)

type SubscriptionInfo struct {
	getter RepositoryGetter
}

func NewSubscriptionInfo(getter RepositoryGetter) *SubscriptionInfo {
	return &SubscriptionInfo{getter: getter}
}

func (u *SubscriptionInfo) Fetch(ctx context.Context) ([]domain.RepositoryInfo, error) {
	return u.getter.GetSubscriptionsInfo(ctx)
}
