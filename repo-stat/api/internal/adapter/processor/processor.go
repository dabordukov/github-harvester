package processor

import (
	"context"
	"log/slog"
	"net/http"

	"repo-stat/api/internal/domain"
	processorpb "repo-stat/proto/processor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   processorpb.ProcessorClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(address)
	if err != nil {
		return nil, err
	}

	return &Client{
		log:  log,
		conn: conn,
		pb:   processorpb.NewProcessorClient(conn),
	}, nil
}

func (c *Client) Ping(ctx context.Context) domain.PingStatus {
	_, err := c.pb.Ping(ctx, &processorpb.PingRequest{})
	if err != nil {
		c.log.Error("Processor ping failed", "error", err)
		return domain.PingStatusDown
	}

	return domain.PingStatusUp
}

func (c *Client) GetRepository(ctx context.Context, owner, repo string) (*domain.RepositoryInfo, error) {
	ctx = appendAuthorizationFromHTTPContext(ctx)

	resp, err := c.pb.GetRepository(ctx, &processorpb.GetRepoRequest{
		Owner:    owner,
		RepoName: repo,
	})
	if err != nil {
		c.log.Error("processor get repository failed", "error", err, "owner", owner, "repo", repo)
		return nil, err
	}

	return &domain.RepositoryInfo{
		Name:         resp.GetName(),
		Owner:        resp.GetOwner(),
		Description:  resp.GetDescription(),
		Forks:        resp.GetForks(),
		Stars:        resp.GetStars(),
		CreatedAt:    resp.GetCreatedAt(),
		CommitsCount: resp.GetCommitsCount(),
	}, nil
}

func appendAuthorizationFromHTTPContext(ctx context.Context) context.Context {
	authHeader := authorizationFromHTTPRequest(ctx)
	if authHeader == "" {
		return ctx
	}

	return metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
}

func authorizationFromHTTPRequest(ctx context.Context) string {
	req, ok := ctx.Value(http.ServerContextKey).(*http.Request)
	if !ok || req == nil {
		return ""
	}

	return req.Header.Get("Authorization")
}

func (c *Client) Close() error {
	return c.conn.Close()
}
