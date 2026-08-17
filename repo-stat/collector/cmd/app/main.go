package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"repo-stat/collector/config"
	collectoradapter "repo-stat/collector/internal/adapter"
	collectorgrpc "repo-stat/collector/internal/controller/grpc"
	"repo-stat/collector/internal/usecase"
	"repo-stat/platform/logger"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)
	log := logger.MustMakeLogger(cfg.Logger.LogLevel)

	log.Info("starting server...")
	log.Debug("debug messages are enabled")

	requestProducer, err := collectoradapter.NewKafkaProducer(cfg.Kafka.Brokers, cfg.Kafka.RequestTopic, log)
	if err != nil {
		log.Warn("kafka request producer unavailable; collector will remain degraded", "error", err)
	}
	resultProducer, err := collectoradapter.NewKafkaProducer(cfg.Kafka.Brokers, cfg.Kafka.ResultTopic, log)
	if err != nil {
		log.Warn("kafka result producer unavailable; collector will remain degraded", "error", err)
	} else {
		defer resultProducer.Close()
	}
	if requestProducer != nil {
		defer requestProducer.Close()
	}

	githubAdapter := collectoradapter.NewGitHubAdapter()
	subscriberClient, err := collectoradapter.NewSubscriberClient(cfg.Services.Subscriber, log)
	if err != nil {
		log.Warn("subscriber client unavailable", "error", err)
	} else {
		defer func() { _ = subscriberClient.Close() }()
	}

	consumer, err := collectoradapter.NewKafkaConsumer(cfg.Kafka.Brokers, cfg.Kafka.RequestTopic, cfg.Kafka.GroupID, log)
	if err != nil {
		log.Warn("kafka consumer unavailable; continuing without async collection", "error", err)
	} else {
		defer consumer.Close()
		go consumer.StartCollectLoop(ctx, func(ctx context.Context, req usecase.CollectRequest) error {
			if resultProducer == nil {
				return fmt.Errorf("kafka result producer not initialized")
			}
			repo, err := githubAdapter.FetchAll(ctx, req.Owner, req.RepoName)
			if err != nil {
				return err
			}
			return resultProducer.PublishCollectResult(ctx, req, repo)
		})
	}

	if err := startSubscriptionRefresh(ctx, cfg, log, subscriberClient, requestProducer); err != nil {
		log.Warn("subscription refresh scheduler unavailable", "error", err)
	}

	srv, cleanup, err := collectorgrpc.NewServerHandler(log, cfg)
	if err != nil {
		log.Error("error creating server", "error", err)
		return err
	}
	defer cleanup()

	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("run grpc server: %w", err)
	}
	return nil
}

func startSubscriptionRefresh(ctx context.Context, cfg config.Config, log *slog.Logger, subscriberClient *collectoradapter.SubscriberClient, requestProducer *collectoradapter.KafkaProducer) error {
	if subscriberClient == nil || requestProducer == nil {
		return fmt.Errorf("subscriber client or kafka producer is not initialized")
	}
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				list, err := subscriberClient.ListSubscriptions(ctx)
				if err != nil {
					log.Error("list subscriptions for async refresh failed", "error", err)
					continue
				}
				for _, subscription := range list {
					if err := requestProducer.PublishCollectRequest(ctx, subscription.GetOwner(), subscription.GetRepoName()); err != nil {
						log.Error("publish refresh request failed", "error", err, "owner", subscription.GetOwner(), "repo", subscription.GetRepoName())
					}
				}
			}
		}
	}()
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
