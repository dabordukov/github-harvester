package usecase

import (
	"context"

	"repo-stat/api/internal/domain"
)

type Subscription struct {
	manager SubscriptionManager
}

func NewSubscription(manager SubscriptionManager) *Subscription {
	return &Subscription{manager: manager}
}

func (u *Subscription) Create(ctx context.Context, owner, repo string) (*domain.Subscription, error) {
	return u.manager.CreateSubscription(ctx, owner, repo)
}

func (u *Subscription) Delete(ctx context.Context, owner, repo string) error {
	return u.manager.DeleteSubscription(ctx, owner, repo)
}

func (u *Subscription) List(ctx context.Context) ([]domain.Subscription, error) {
	return u.manager.ListSubscriptions(ctx)
}
