CREATE TABLE IF NOT EXISTS
    repositories (
        id BIGSERIAL PRIMARY KEY,
        repo_owner TEXT NOT NULL,
        repo_name TEXT NOT NULL,
        description TEXT,
        forks INT NOT NULL DEFAULT 0,
        stars INT NOT NULL DEFAULT 0,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        commits_count INT NOT NULL DEFAULT 0,
        CONSTRAINT repositories_owner_repo_name_key UNIQUE (repo_owner, repo_name)
    );

CREATE INDEX IF NOT EXISTS repositories_owner_name_idx ON repositories (repo_owner, repo_name);
