package http

import (
	"log/slog"
	"net/http"

	"repo-stat/api/internal/usecase"

	httpSwagger "github.com/swaggo/http-swagger"
)

func AddRoutes(
	mux *http.ServeMux,
	log *slog.Logger,
	ping *usecase.Ping,
	repositoryInfo *usecase.RepositoryInfo,
	subscription *usecase.Subscription,
	subscriptionInfo *usecase.SubscriptionInfo,
) {
	mux.Handle("GET /api/ping", NewPingHandler(log, ping))
	mux.Handle("GET /api/repositories/info", NewRepositoryInfoHandler(log, repositoryInfo))
	mux.Handle("GET /api/subscriptions", NewListSubscriptionsHandler(log, subscription))
	mux.Handle("POST /api/subscriptions", NewCreateSubscriptionHandler(log, subscription))
	mux.Handle("DELETE /api/subscriptions/{owner}/{repo}", NewDeleteSubscriptionHandler(log, subscription))
	mux.Handle("GET /api/subscriptions/info", NewSubscriptionsInfoHandler(log, subscriptionInfo))
	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)
}
