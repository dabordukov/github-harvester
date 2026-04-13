package adapter

import (
	"context"

	"github-harvester/internal/gateway/model"
	"github-harvester/internal/gateway/service"
	"github-harvester/internal/pkg/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CollectorAdapter struct {
	client pb.CollectorClient
}

func NewCollectorAdapter(client pb.CollectorClient) *CollectorAdapter {
	return &CollectorAdapter{
		client: client,
	}
}

func (ca *CollectorAdapter) GetRepo(ctx context.Context, owner, repoName string) (*model.RepositoryModel, error) {
	res, err := ca.client.GetRepository(ctx, &pb.GetRepoRequest{
		Owner:    owner,
		RepoName: repoName,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, service.ErrNotFound
			case codes.Unauthenticated:
				return nil, service.ErrUnauthorized
			case codes.ResourceExhausted:
				return nil, service.ErrRateLimited
			}
		}
		return nil, err
	}

	return &model.RepositoryModel{
		Name:         res.Name,
		Owner:        res.Owner,
		Description:  res.Description,
		Forks:        res.Forks,
		Stars:        res.Stars,
		CreatedAt:    res.CreatedAt,
		CommitsCount: res.CommitsCount,
	}, nil
}
