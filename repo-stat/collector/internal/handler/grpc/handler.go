package grpc

import (
	"log/slog"

	"repo-stat/collector/config"
	"repo-stat/collector/internal/adapter"
	"repo-stat/collector/internal/service"
	grpcserver "repo-stat/platform/grpcserver"
	collectorpb "repo-stat/proto/collector"
)

func NewServer(log *slog.Logger, cfg config.Config) (*grpcserver.Server, error) {
	githubAdapter := adapter.NewGitHubAdapter()
	collectorService := service.NewCollectorService(githubAdapter)
	collectorServer := NewHandler(log, collectorService)

	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return nil, err
	}

	collectorpb.RegisterCollectorServer(srv.GRPC(), collectorServer)

	return srv, nil
}
