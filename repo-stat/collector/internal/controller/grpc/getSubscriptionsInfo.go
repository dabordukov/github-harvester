package grpc

import (
	"context"
	"errors"

	service "repo-stat/collector/internal/usecase"
	collectorpb "repo-stat/proto/collector"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetSubscriptionsInfo(
	ctx context.Context,
	_ *collectorpb.GetSubscriptionsInfoRequest,
) (*collectorpb.GetSubscriptionsInfoResponse, error) {
	repositories, err := h.service.GetSubscriptionsData(ctx)
	if err != nil {
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

	response := &collectorpb.GetSubscriptionsInfoResponse{
		Repositories: make([]*collectorpb.GetRepoResponse, 0, len(repositories)),
	}
	for _, repository := range repositories {
		response.Repositories = append(response.Repositories, &collectorpb.GetRepoResponse{
			FullName:     repository.FullName,
			Name:         repository.Name,
			Owner:        repository.Owner,
			Description:  repository.Description,
			Forks:        repository.Forks,
			Stars:        repository.Stars,
			CreatedAt:    repository.CreatedAt,
			CommitsCount: repository.CommitsCount,
		})
	}

	return response, nil
}
