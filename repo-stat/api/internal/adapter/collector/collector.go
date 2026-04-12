package collector

import (
	"log/slog"

	collectorpb "repo-stat/proto/collector"

	"google.golang.org/grpc"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   collectorpb.CollectorClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(address)
	if err != nil {
		return nil, err
	}

	return &Client{
		log:  log,
		conn: conn,
		pb:   collectorpb.NewCollectorClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
