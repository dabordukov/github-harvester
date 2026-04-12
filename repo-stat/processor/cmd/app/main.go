package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	"repo-stat/processor/config"
	collectoradapter "repo-stat/processor/internal/adapter/collector"
	grpccontroller "repo-stat/processor/internal/controller/grpc"
	service "repo-stat/processor/internal/usecase"
	processorpb "repo-stat/proto/processor"
)

func run(ctx context.Context) error {
	// config
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	// logger

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)

	log.Info("starting server...")
	log.Debug("debug messages are enabled")

	collectorClient, err := collectoradapter.NewClient(cfg.Services.Collector, log)
	if err != nil {
		log.Error("cannot init collector adapter", "error", err)
		return err
	}
	defer func() {
		if err := collectorClient.Close(); err != nil {
			log.Error("cannot close collector client", "error", err)
		}
	}()

	processorService := service.NewProcessorService(collectorClient)
	processorServer := grpccontroller.NewServer(log, processorService)

	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	processorpb.RegisterProcessorServer(srv.GRPC(), processorServer)

	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("run grpc server: %w", err)
	}
	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	if err := run(ctx); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			fmt.Printf("launching server error: %s\n", err)
		}
		cancel()
		os.Exit(1)
	}
	cancel()
}
