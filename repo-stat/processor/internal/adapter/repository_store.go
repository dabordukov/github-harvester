package adapter

import (
	"context"
	"fmt"
	"time"

	"repo-stat/processor/internal/domain"
	"repo-stat/processor/internal/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type RepositoryStore struct {
	queries *sqlc.Queries
}

func NewRepositoryStore(queries *sqlc.Queries) *RepositoryStore {
	return &RepositoryStore{queries: queries}
}

func (s *RepositoryStore) Get(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	if s == nil || s.queries == nil {
		return nil, fmt.Errorf("repository store is not initialized")
	}

	row, err := s.queries.GetRepository(ctx, sqlc.GetRepositoryParams{
		RepoOwner: owner,
		RepoName:  repo,
	})
	if err != nil {
		return nil, err
	}

	return &domain.Repository{
		Owner:        row.RepoOwner,
		Name:         row.RepoName,
		FullName:     fmt.Sprintf("%s/%s", row.RepoOwner, row.RepoName),
		Forks:        int64(row.Forks),
		Stars:        int64(row.Stars),
		CreatedAt:    row.CreatedAt.Time.Format(time.RFC3339),
		CommitsCount: 0,
	}, nil
}

func (s *RepositoryStore) Save(ctx context.Context, repository *domain.Repository) error {
	if s == nil || s.queries == nil {
		return fmt.Errorf("repository store is not initialized")
	}
	if repository == nil {
		return fmt.Errorf("repository is nil")
	}

	createdAt := parseCreatedAt(repository.CreatedAt)
	_, err := s.queries.UpsertRepository(ctx, sqlc.UpsertRepositoryParams{
		RepoOwner:    repository.Owner,
		RepoName:     repository.Name,
		Description:  pgtype.Text{String: repository.Description, Valid: true},
		Forks:        int32(repository.Forks),
		Stars:        int32(repository.Stars),
		CreatedAt:    pgtype.Timestamptz{Time: createdAt, Valid: true},
		CommitsCount: int32(repository.CommitsCount),
	})
	return err
}

func parseCreatedAt(value string) time.Time {
	if value == "" {
		return time.Now().UTC()
	}

	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Now().UTC()
	}
	return parsed.UTC()
}
