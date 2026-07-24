package grpc

import (
	"context"

	subscriberpb "repo-stat/proto/subscriber"
)

func (s *Server) Ping(ctx context.Context, _ *subscriberpb.PingRequest) (*subscriberpb.PingResponse, error) {
	s.log.Debug("subscriberp ping request received")

	return &subscriberpb.PingResponse{
		Reply: s.ping.Execute(ctx),
	}, nil
}
