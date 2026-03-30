package collector

import (
	"context"
	"log/slog"

	"repo-stat/processor/internal/domain"
	collectorpb "repo-stat/proto/collector"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   collectorpb.CollectorClient
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
		pb:   collectorpb.NewCollectorClient(conn),
	}, nil
}

func (c *Client) GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	ctx = forwardAuthorization(ctx)

	resp, err := c.pb.GetRepository(ctx, &collectorpb.GetRepoRequest{
		Owner:    owner,
		RepoName: repo,
	})
	if err != nil {
		c.log.Error("collector get repository failed", "error", err, "owner", owner, "repo", repo)
		return nil, err
	}

	return &domain.Repository{
		FullName:     resp.GetFullName(),
		Name:         resp.GetName(),
		Owner:        resp.GetOwner(),
		Description:  resp.GetDescription(),
		Forks:        resp.GetForks(),
		Stars:        resp.GetStars(),
		CreatedAt:    resp.GetCreatedAt(),
		CommitsCount: resp.GetCommitsCount(),
	}, nil
}

func (c *Client) Ping(ctx context.Context) string {
	ctx = forwardAuthorization(ctx)

	_, err := c.pb.Ping(ctx, &collectorpb.PingRequest{})
	if err != nil {
		c.log.Error("collector ping failed", "error", err)
		return "down"
	}

	return "pong"
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func forwardAuthorization(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}

	tokens := md.Get("authorization")
	if len(tokens) == 0 || tokens[0] == "" {
		return ctx
	}

	return metadata.AppendToOutgoingContext(ctx, "authorization", tokens[0])
}
