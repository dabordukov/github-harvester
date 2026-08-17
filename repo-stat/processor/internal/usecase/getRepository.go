package usecase

import (
	"context"
	"fmt"

	"repo-stat/processor/internal/domain"
)

type ProcessorService struct {
	collector CollectorProvider
	store     RepositoryStore
	publisher KafkaPublisher
}

type ProcessorDependency struct {
	Store     RepositoryStore
	Publisher KafkaPublisher
}

func NewProcessorService(collector CollectorProvider, deps ...ProcessorDependency) *ProcessorService {
	service := &ProcessorService{collector: collector}
	if len(deps) > 0 {
		service.store = deps[0].Store
		service.publisher = deps[0].Publisher
	}
	return service
}

func (s *ProcessorService) GetRepositoryData(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	if s.store != nil {
		if cached, err := s.store.Get(ctx, owner, repo); err == nil && cached != nil {
			return cached, nil
		}
	}

	if s.publisher != nil {
		if err := s.publisher.PublishCollectRequest(ctx, owner, repo); err != nil {
			return nil, err
		}
	}

	if s.store != nil {
		if cached, err := s.store.Get(ctx, owner, repo); err == nil && cached != nil {
			return cached, nil
		}
	}

	if s.collector != nil {
		return s.collector.GetRepository(ctx, owner, repo)
	}

	return nil, fmt.Errorf("repository %s/%s not found in cache and no collector available", owner, repo)
}

func (s *ProcessorService) GetSubscriptionsInfo(ctx context.Context) ([]domain.Repository, error) {
	if s.collector != nil {
		return s.collector.GetSubscriptionsInfo(ctx)
	}
	return nil, nil
}
