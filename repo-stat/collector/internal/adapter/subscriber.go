package adapter

import (
	"context"
	"log/slog"

	subscriberpb "repo-stat/proto/subscriber"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SubscriberClient struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   subscriberpb.SubscriberClient
}

func NewSubscriberClient(address string, log *slog.Logger) (*SubscriberClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &SubscriberClient{
		log:  log,
		conn: conn,
		pb:   subscriberpb.NewSubscriberClient(conn),
	}, nil
}

func (c *SubscriberClient) ListSubscriptions(ctx context.Context) ([]*subscriberpb.Subscription, error) {
	resp, err := c.pb.ListSubscriptions(ctx, &subscriberpb.ListSubscriptionsRequest{})
	if err != nil {
		c.log.Error("subscriber list subscriptions failed", "error", err)
		return nil, err
	}

	return resp.GetSubscriptions(), nil
}

func (c *SubscriberClient) Close() error {
	return c.conn.Close()
}
