-- name: CreateSubscription :one
INSERT INTO subscriptions (
    repo_owner,
    repo_name
) VALUES (
    $1,
    $2
)
RETURNING id, repo_owner, repo_name, created_at;

-- name: DeleteSubscription :execrows
DELETE FROM subscriptions
WHERE repo_owner = $1 AND repo_name = $2;

-- name: ListSubscriptions :many
SELECT id, repo_owner, repo_name, created_at
FROM subscriptions
ORDER BY created_at DESC, id DESC;
