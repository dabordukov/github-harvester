package grpc

import (
	"context"

	"repo-stat/proto/collector"
)

func (h *Handler) Ping(ctx context.Context, _ *collector.PingRequest) (*collector.PingResponse, error) {
	h.log.Debug("collector ping request received")

	return &collector.PingResponse{
		Reply: h.ping.Execute(ctx),
	}, nil
}
