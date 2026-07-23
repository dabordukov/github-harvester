package grpc

import (
	"context"
	"errors"
	"log"
	"log/slog"

	service "repo-stat/collector/internal/usecase"
	"repo-stat/proto/collector"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServiceInterface interface {
	GetRepositoryData(ctx context.Context, owner, repo string) (*service.RepositoryModel, error)
	GetSubscriptionsData(ctx context.Context) ([]service.RepositoryModel, error)
}

type Handler struct {
	collector.UnimplementedCollectorServer
	log     *slog.Logger
	service ServiceInterface
	ping    *service.Ping
}

func NewHandler(log *slog.Logger, service ServiceInterface, ping *service.Ping) *Handler {
	return &Handler{
		log:     log,
		service: service,
		ping:    ping,
	}
}

func (h *Handler) GetRepository(ctx context.Context, req *collector.GetRepoRequest) (*collector.GetRepoResponse, error) {
	if req.GetOwner() == "" {
		return nil, status.Error(codes.InvalidArgument, "owner is required")
	}

	if req.GetRepoName() == "" {
		return nil, status.Error(codes.InvalidArgument, "repoName is required")
	}

	res, err := h.service.GetRepositoryData(ctx, req.GetOwner(), req.GetRepoName())
	if err != nil {
		log.Printf("Error: %v", err)

		switch {
		case errors.Is(err, service.ErrNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, service.ErrRateLimited):
			return nil, status.Error(codes.ResourceExhausted, err.Error())
		case errors.Is(err, service.ErrUnauthorized):
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &collector.GetRepoResponse{
		FullName:     res.FullName,
		Name:         res.Name,
		Owner:        res.Owner,
		Description:  res.Description,
		Forks:        res.Forks,
		Stars:        res.Stars,
		CreatedAt:    res.CreatedAt,
		CommitsCount: res.CommitsCount,
	}, nil
}
