package usecase

import "errors"

var (
	ErrRepoNotFound     = errors.New("repository not found in subscriptions")
	ErrAlreadyExists    = errors.New("subscription already exists")
	ErrOwnerRequired    = errors.New("owner is required")
	ErrRepoNameRequired = errors.New("repo_name is required")
)
