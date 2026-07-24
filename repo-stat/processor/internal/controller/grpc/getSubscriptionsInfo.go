package grpc

import (
	"context"

	processorpb "repo-stat/proto/processor"

	"google.golang.org/grpc/status"
)

func (s *Server) GetSubscriptionsInfo(
	ctx context.Context,
	_ *processorpb.GetSubscriptionsInfoRequest,
) (*processorpb.GetSubscriptionsInfoResponse, error) {
	repositories, err := s.service.GetSubscriptionsInfo(ctx)
	if err != nil {
		s.log.Error("processor get subscriptions info failed", "error", err)
		return nil, status.Convert(err).Err()
	}

	response := &processorpb.GetSubscriptionsInfoResponse{
		Repositories: make([]*processorpb.GetRepoResponse, 0, len(repositories)),
	}
	for _, repo := range repositories {
		response.Repositories = append(response.Repositories, &processorpb.GetRepoResponse{
			FullName:     repo.FullName,
			Name:         repo.Name,
			Owner:        repo.Owner,
			Description:  repo.Description,
			Forks:        repo.Forks,
			Stars:        repo.Stars,
			CreatedAt:    repo.CreatedAt,
			CommitsCount: repo.CommitsCount,
		})
	}

	return response, nil
}
