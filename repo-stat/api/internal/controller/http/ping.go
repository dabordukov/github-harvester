package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

func NewPingHandler(log *slog.Logger, ping *usecase.Ping) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result := ping.Execute(r.Context())

		services := make([]dto.PingServiceInfo, 0, len(result.Services))
		for _, service := range result.Services {
			services = append(services, dto.PingServiceInfo{
				Name:   service.Name,
				Status: string(service.Status),
			})
		}

		response := dto.PingResponse{
			Status:   result.Status,
			Services: services,
		}

		w.Header().Set("Content-Type", "application/json")
		statusCode := http.StatusOK
		if result.Status != "ok" {
			statusCode = http.StatusServiceUnavailable
		}
		w.WriteHeader(statusCode)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to write ping response", "error", err)
		}
	}
}
