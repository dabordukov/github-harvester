package grpc

import (
	"context"

	"repo-stat/proto/processor"
)

func (s *Server) Ping(ctx context.Context, _ *processor.PingRequest) (*processor.PingResponse, error) {
	s.log.Debug("processor ping request received")

	return &processor.PingResponse{
		Reply: s.service.Ping(ctx),
	}, nil
}
