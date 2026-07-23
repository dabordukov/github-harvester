package grpc

import (
	"context"
	"errors"
	"log"
	"log/slog"

	"repo-stat/collector/internal/usecase"
	collectorpb "repo-stat/proto/collector"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServiceInterface interface {
	GetRepositoryData(ctx context.Context, owner, repo string) (*usecase.RepositoryModel, error)
	GetSubscriptionsData(ctx context.Context) ([]usecase.RepositoryModel, error)
}

type Handler struct {
	collectorpb.UnimplementedCollectorServer
	log     *slog.Logger
	service ServiceInterface
	ping    *usecase.Ping
}

func NewHandler(log *slog.Logger, service ServiceInterface, ping *usecase.Ping) *Handler {
	return &Handler{
		log:     log,
		service: service,
		ping:    ping,
	}
}

func (h *Handler) GetRepository(ctx context.Context, req *collectorpb.GetRepoRequest) (*collectorpb.GetRepoResponse, error) {
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
		case errors.Is(err, usecase.ErrNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, usecase.ErrRateLimited):
			return nil, status.Error(codes.ResourceExhausted, err.Error())
		case errors.Is(err, usecase.ErrUnauthorized):
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &collectorpb.GetRepoResponse{
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
