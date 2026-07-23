package grpc

import (
	"log/slog"

	"repo-stat/collector/config"
	"repo-stat/collector/internal/adapter"
	service "repo-stat/collector/internal/usecase"
	grpcserver "repo-stat/platform/grpcserver"
	collectorpb "repo-stat/proto/collector"
)

func NewServerHandler(log *slog.Logger, cfg config.Config) (*grpcserver.Server, error) {
	githubAdapter := adapter.NewGitHubAdapter()
	subscriberClient, err := adapter.NewSubscriberClient(cfg.Services.Subscriber, log)
	if err != nil {
		return nil, err
	}
	subscriberClient.Close()

	collectorService := service.NewCollectorService(githubAdapter, subscriberClient)
	pingService := service.NewPing()
	collectorServer := NewHandler(log, collectorService, pingService)

	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return nil, err
	}

	collectorpb.RegisterCollectorServer(srv.GRPC(), collectorServer)

	return srv, nil
}
