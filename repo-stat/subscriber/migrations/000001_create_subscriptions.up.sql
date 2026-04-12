CREATE TABLE IF NOT EXISTS
    subscriptions (
        id BIGSERIAL PRIMARY KEY,
        repo_owner TEXT NOT NULL,
        repo_name TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        CONSTRAINT subscriptions_owner_repo_name_key UNIQUE (repo_owner, repo_name)
    );

CREATE INDEX subscriptions_created_at_id_desc_idx ON subscriptions (created_at DESC, id DESC);
