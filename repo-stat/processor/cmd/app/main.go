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
	kafkaadapter "repo-stat/processor/internal/adapter"
	collectoradapter "repo-stat/processor/internal/adapter/collector"
	grpccontroller "repo-stat/processor/internal/controller/grpc"
	"repo-stat/processor/internal/sqlc"
	"repo-stat/processor/internal/usecase"
	processorpb "repo-stat/proto/processor"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)
	log := logger.MustMakeLogger(cfg.Logger.LogLevel)

	log.Info("starting server...")
	log.Debug("debug messages are enabled")

	pool, err := pgxpool.New(ctx, cfg.Database.DSN)
	if err != nil {
		return fmt.Errorf("create pgx pool: %w", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}
	if err := runMigrations(cfg); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	queries := sqlc.New(pool)
	repositoryStore := kafkaadapter.NewRepositoryStore(queries)
	producer, err := kafkaadapter.NewKafkaProducer(cfg.Kafka.Brokers, cfg.Kafka.RequestTopic, log)
	if err != nil {
		log.Warn("kafka producer unavailable; continuing without async request publishing", "error", err)
	} else {
		defer producer.Close()
	}

	consumer, err := kafkaadapter.NewKafkaConsumer(cfg.Kafka.Brokers, cfg.Kafka.ResultTopic, cfg.Kafka.GroupID, log)
	if err != nil {
		log.Warn("kafka consumer unavailable; continuing without async result processing", "error", err)
	} else {
		defer consumer.Close()
		if err := consumer.StartResultConsumer(ctx, repositoryStore); err != nil {
			log.Warn("cannot start kafka result consumer", "error", err)
		}
	}

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

	processorService := usecase.NewProcessorService(collectorClient, usecase.ProcessorDependency{
		Store:     repositoryStore,
		Publisher: producer,
	})
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

func runMigrations(cfg config.Config) error {
	migrator, err := migrate.New(cfg.Database.MigrationsPath, cfg.Database.DSN)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = migrator.Close()
	}()

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return err
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
