package http

import (
	"context"
	"log/slog"
	"net/http"

	"repo-stat/api/config"
	"repo-stat/api/internal/adapter/processor"
	"repo-stat/api/internal/adapter/subscriber"
	"repo-stat/api/internal/usecase"
)

func NewHandler(_ context.Context, log *slog.Logger, cfg config.Config) (http.Handler, func(), error) {
	processorClient, err := processor.NewClient(cfg.Services.Processor, log)
	if err != nil {
		log.Error("cannot init processor adapter", "error", err)
		return nil, nil, err
	}

	subscriberClient, err := subscriber.NewClient(cfg.Services.Subscriber, log)
	if err != nil {
		log.Error("cannot init subscriber adapter", "error", err)
		if err := processorClient.Close(); err != nil {
			log.Error("failed to close processor client", "error", err)
		}
		return nil, nil, err
	}

	pingUseCase := usecase.NewPing(processorClient, subscriberClient)
	repositoryInfoUseCase := usecase.NewRepositoryInfo(processorClient)
	subscriptionUseCase := usecase.NewSubscription(subscriberClient)
	subscriptionInfoUseCase := usecase.NewSubscriptionInfo(processorClient)

	handler := http.NewServeMux()
	AddRoutes(handler, log, pingUseCase, repositoryInfoUseCase, subscriptionUseCase, subscriptionInfoUseCase)

	cleanup := func() {
		log.Info("closing gRPC clients...")
		if err := processorClient.Close(); err != nil {
			log.Error("failed to close processor client", "error", err)
		}
		if err := subscriberClient.Close(); err != nil {
			log.Error("failed to close subscriber client", "error", err)
		}
	}

	return handler, cleanup, nil
}
