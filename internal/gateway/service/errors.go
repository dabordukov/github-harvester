package service

import "errors"

var (
	ErrInvalidInput = errors.New("owner and repository name are required")
	ErrNotFound     = errors.New("repository not found")
	ErrUnauthorized = errors.New("invalid or expired token")
	ErrRateLimited  = errors.New("too many requests to github")
)
