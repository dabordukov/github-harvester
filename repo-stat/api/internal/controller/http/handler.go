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

func NewHandler(ctx context.Context, log *slog.Logger, cfg config.Config) (http.Handler, error) {
	processorClient, err := processor.NewClient(cfg.Services.Processor, log)
	if err != nil {
		log.Error("cannot init processor adapter", "error", err)
		return nil, err
	}

	subscriberClient, err := subscriber.NewClient(cfg.Services.Subscriber, log)
	if err != nil {
		log.Error("cannot init subscriber adapter", "error", err)
		return nil, err
	}

	pingUseCase := usecase.NewPing(processorClient, subscriberClient)
	repositoryInfoUseCase := usecase.NewRepositoryInfo(processorClient)

	mux := http.NewServeMux()
	AddRoutes(mux, log, pingUseCase, repositoryInfoUseCase)

	var handler http.Handler = mux
	_ = ctx
	return handler, nil
}
