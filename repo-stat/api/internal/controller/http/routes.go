package http

import (
	"log/slog"
	"net/http"

	"repo-stat/api/internal/usecase"

	httpSwagger "github.com/swaggo/http-swagger"
)

func AddRoutes(mux *http.ServeMux, log *slog.Logger, ping *usecase.Ping, repositoryInfo *usecase.RepositoryInfo) {
	mux.Handle("GET /api/ping", NewPingHandler(log, ping))
	mux.Handle("GET /api/repositories/info", NewRepositoryInfoHandler(log, repositoryInfo))
	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)
}
