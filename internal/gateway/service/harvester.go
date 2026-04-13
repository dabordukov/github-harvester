package service

import (
	"context"
	"fmt"

	"github-harvester/internal/gateway/model"
)

type CollectorProvider interface {
	GetRepo(ctx context.Context, owner, repo string) (*model.RepositoryModel, error)
}

type HarvesterService struct {
	collector CollectorProvider
}

func NewHarvesterService(cp CollectorProvider) *HarvesterService {
	return &HarvesterService{
		collector: cp,
	}
}

func (h *HarvesterService) GetRepositoryInfo(ctx context.Context, owner, repoName string) (*model.RepositoryModel, error) {
	if owner == "" || repoName == "" {
		return nil, fmt.Errorf("repoName is required")
	}

	res, err := h.collector.GetRepo(ctx, owner, repoName)
	if err != nil {
		return nil, err
	}

	return res, nil
}
