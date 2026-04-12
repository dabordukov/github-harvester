package service

import (
	"context"

	"repo-stat/processor/internal/domain"
)

type ProcessorService struct {
	collector CollectorProvider
}

func NewProcessorService(collector CollectorProvider) *ProcessorService {
	return &ProcessorService{collector: collector}
}

func (s *ProcessorService) GetRepositoryData(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	return s.collector.GetRepository(ctx, owner, repo)
}

func (s *ProcessorService) GetSubscriptionsInfo(ctx context.Context) ([]domain.Repository, error) {
	return s.collector.GetSubscriptionsInfo(ctx)
}
