-- name: GetRepository :one
SELECT
    id,
    repo_owner,
    repo_name,
    forks,
    stars,
    created_at
FROM
    repositories
WHERE
    repo_owner = $1
    AND repo_name = $2;

-- name: UpsertRepository :one
INSERT INTO
    repositories (
        repo_owner,
        repo_name,
        description,
        forks,
        stars,
        created_at,
        commits_count
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (repo_owner, repo_name)
DO
UPDATE
SET
    forks = EXCLUDED.forks,
    stars = EXCLUDED.stars,
    created_at = EXCLUDED.created_at
RETURNING
    id,
    repo_owner,
    repo_name,
    forks,
    stars,
    created_at;
