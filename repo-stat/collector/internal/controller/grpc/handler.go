package grpc

import (
	"log/slog"

	"repo-stat/collector/config"
	"repo-stat/collector/internal/adapter"
	"repo-stat/collector/internal/usecase"
	grpcserver "repo-stat/platform/grpcserver"
	collectorpb "repo-stat/proto/collector"
)

func NewServerHandler(log *slog.Logger, cfg config.Config) (*grpcserver.Server, func(), error) {
	githubAdapter := adapter.NewGitHubAdapter()
	subscriberClient, err := adapter.NewSubscriberClient(cfg.Services.Subscriber, log)
	if err != nil {
		return nil, nil, err
	}

	collectorService := usecase.NewCollectorService(githubAdapter, subscriberClient, log)
	pingService := usecase.NewPing()
	collectorServer := NewHandler(log, collectorService, pingService)

	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		if err := subscriberClient.Close(); err != nil {
			log.Error("failed to close subscriber client", "error", err)
		}
		return nil, nil, err
	}

	collectorpb.RegisterCollectorServer(srv.GRPC(), collectorServer)

	cleanup := func() {
		log.Info("closing gRPC clients...")
		if err := subscriberClient.Close(); err != nil {
			log.Error("failed to close subscriber client", "error", err)
		}
	}

	return srv, cleanup, nil
}
