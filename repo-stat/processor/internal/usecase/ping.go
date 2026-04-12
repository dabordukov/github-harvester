package service

import "context"

func (s *ProcessorService) Ping(context.Context) string {
	return "pong"
}
