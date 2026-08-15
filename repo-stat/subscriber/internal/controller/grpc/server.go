package grpc

import (
	"log/slog"

	subscriberpb "repo-stat/proto/subscriber"
	"repo-stat/subscriber/internal/usecase"
)

type Server struct {
	subscriberpb.UnimplementedSubscriberServer
	log          *slog.Logger
	ping         *usecase.Ping
	subscription *usecase.Subscription
}

func NewServer(log *slog.Logger, ping *usecase.Ping, subscription *usecase.Subscription) *Server {
	return &Server{
		log:          log,
		ping:         ping,
		subscription: subscription,
	}
}
