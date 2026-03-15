package adapter

import "errors"

var (
	ErrNotFound     = errors.New("provider: resource not found")
	ErrUnauthorized = errors.New("provider: unauthorized")
	ErrRateLimited  = errors.New("provider: rate limit exceeded")
)
