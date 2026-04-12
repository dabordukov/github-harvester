package grpc

import (
	"context"
	"log/slog"

	"repo-stat/processor/internal/domain"
	"repo-stat/proto/processor"
)

type ServiceInterface interface {
	GetRepositoryData(ctx context.Context, owner, repo string) (*domain.Repository, error)
	Ping(ctx context.Context) string
}

type Server struct {
	processor.UnimplementedProcessorServer
	log     *slog.Logger
	service ServiceInterface
}

func NewServer(log *slog.Logger, service ServiceInterface) *Server {
	return &Server{
		log:     log,
		service: service,
	}
}
