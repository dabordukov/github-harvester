package handlers

import (
	"context"
	"errors"
	"log"

	"github-harvester/internal/collector/adapter"
	"github-harvester/internal/collector/service"
	"github-harvester/internal/pkg/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServiceInterface interface {
	GetRepositoryData(ctx context.Context, owner, repo string) (*service.RepositoryModel, error)
}

type Handler struct {
	pb.UnimplementedCollectorServer
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{
		UnimplementedCollectorServer: pb.UnimplementedCollectorServer{},
		service:                      service,
	}
}

func (h *Handler) GetRepository(ctx context.Context, req *pb.GetRepoRequest) (*pb.GetRepoResponse, error) {
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
		case errors.Is(err, adapter.ErrNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, adapter.ErrRateLimited):
			return nil, status.Error(codes.ResourceExhausted, err.Error())
		case errors.Is(err, adapter.ErrUnauthorized):
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.GetRepoResponse{
		Name:         res.Name,
		Owner:        res.Owner,
		Description:  res.Description,
		Forks:        res.Forks,
		Stars:        res.Stars,
		CreatedAt:    res.CreatedAt,
		CommitsCount: res.CommitsCount,
	}, nil
}
