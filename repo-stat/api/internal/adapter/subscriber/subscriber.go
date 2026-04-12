package subscriber

import (
	"context"
	"log/slog"
	"net/http"

	"repo-stat/api/internal/domain"

	subscriberpb "repo-stat/proto/subscriber"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   subscriberpb.SubscriberClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		log:  log,
		conn: conn,
		pb:   subscriberpb.NewSubscriberClient(conn),
	}, nil
}

func (c *Client) Ping(ctx context.Context) domain.PingStatus {
	_, err := c.pb.Ping(ctx, &subscriberpb.PingRequest{})
	if err != nil {
		c.log.Error("subscriber ping failed", "error", err)
		return domain.PingStatusDown
	}

	return domain.PingStatusUp
}

func (c *Client) CreateSubscription(ctx context.Context, owner, repo string) (*domain.Subscription, error) {
	ctx = appendAuthorizationFromHTTPContext(ctx)

	resp, err := c.pb.CreateSubscription(ctx, &subscriberpb.CreateSubscriptionRequest{
		Subscription: &subscriberpb.Subscription{
			Owner:    owner,
			RepoName: repo,
		},
	})
	if err != nil {
		c.log.Error("subscriber create subscription failed", "error", err, "owner", owner, "repo", repo)
		return nil, err
	}

	subscription := resp.GetSubscription()
	if subscription == nil {
		return &domain.Subscription{Owner: owner, RepoName: repo}, nil
	}

	return &domain.Subscription{
		Owner:    subscription.GetOwner(),
		RepoName: subscription.GetRepoName(),
	}, nil
}

func (c *Client) DeleteSubscription(ctx context.Context, owner, repo string) error {
	ctx = appendAuthorizationFromHTTPContext(ctx)

	_, err := c.pb.DeleteSubscription(ctx, &subscriberpb.DeleteSubscriptionRequest{
		Subscription: &subscriberpb.Subscription{
			Owner:    owner,
			RepoName: repo,
		},
	})
	if err != nil {
		c.log.Error("subscriber delete subscription failed", "error", err, "owner", owner, "repo", repo)
		return err
	}

	return nil
}

func (c *Client) ListSubscriptions(ctx context.Context) ([]domain.Subscription, error) {
	ctx = appendAuthorizationFromHTTPContext(ctx)

	resp, err := c.pb.ListSubscriptions(ctx, &subscriberpb.ListSubscriptionsRequest{})
	if err != nil {
		c.log.Error("subscriber list subscriptions failed", "error", err)
		return nil, err
	}

	subscriptions := make([]domain.Subscription, 0, len(resp.GetSubscriptions()))
	for _, subscription := range resp.GetSubscriptions() {
		subscriptions = append(subscriptions, domain.Subscription{
			Owner:    subscription.GetOwner(),
			RepoName: subscription.GetRepoName(),
		})
	}

	return subscriptions, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func appendAuthorizationFromHTTPContext(ctx context.Context) context.Context {
	req, ok := ctx.Value(http.ServerContextKey).(*http.Request)
	if !ok || req == nil {
		return ctx
	}

	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return ctx
	}

	return metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
}
