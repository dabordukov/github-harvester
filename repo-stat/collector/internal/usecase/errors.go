package service

import "errors"

var (
	ErrNotFound     = errors.New("repository not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrRateLimited  = errors.New("rate limit exceeded")
	ErrUnexpected   = errors.New("unexpected status")
)
