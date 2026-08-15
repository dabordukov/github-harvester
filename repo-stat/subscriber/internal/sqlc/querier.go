package sqlc

import "context"

type Querier interface {
	CreateSubscription(context.Context, CreateSubscriptionParams) (Subscription, error)
	DeleteSubscription(context.Context, DeleteSubscriptionParams) (int64, error)
	ListSubscriptions(context.Context) ([]Subscription, error)
}
