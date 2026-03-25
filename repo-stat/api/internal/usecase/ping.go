package usecase

import (
	"context"

	"repo-stat/api/internal/domain"
)

type Pinger interface {
	Ping(ctx context.Context) domain.PingStatus
}

type Ping struct {
	processor  Pinger
	subscriber Pinger
}

func NewPing(processor, subscriber Pinger) *Ping {
	return &Ping{
		processor:  processor,
		subscriber: subscriber,
	}
}

func (u *Ping) Execute(ctx context.Context) domain.PingResult {
	result := domain.PingResult{
		Status: "ok",
		Services: []domain.PingService{
			{Name: "processor", Status: u.processor.Ping(ctx)},
			{Name: "subscriber", Status: u.subscriber.Ping(ctx)},
		},
	}

	for _, service := range result.Services {
		if service.Status != domain.PingStatusUp {
			result.Status = "degraded"
			break
		}
	}

	return result
}
