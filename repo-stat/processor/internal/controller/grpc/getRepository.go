package grpc

import (
	"context"

	processorpb "repo-stat/proto/processor"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetRepository(ctx context.Context, req *processorpb.GetRepoRequest) (*processorpb.GetRepoResponse, error) {
	if req.GetOwner() == "" {
		return nil, status.Error(codes.InvalidArgument, "owner is required")
	}

	if req.GetRepoName() == "" {
		return nil, status.Error(codes.InvalidArgument, "repo_name is required")
	}

	repo, err := s.service.GetRepositoryData(ctx, req.GetOwner(), req.GetRepoName())
	if err != nil {
		s.log.Error("processor get repository failed", "error", err, "owner", req.GetOwner(), "repo", req.GetRepoName())
		return nil, status.Convert(err).Err()
	}

	return &processorpb.GetRepoResponse{
		FullName:     repo.FullName,
		Name:         repo.Name,
		Owner:        repo.Owner,
		Description:  repo.Description,
		Forks:        repo.Forks,
		Stars:        repo.Stars,
		CreatedAt:    repo.CreatedAt,
		CommitsCount: repo.CommitsCount,
	}, nil
}
