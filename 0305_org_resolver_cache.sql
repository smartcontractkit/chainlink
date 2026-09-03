-- +goose Up

CREATE TABLE IF NOT EXISTS cre.org_resolver_cache (
    workflow_owner TEXT        NOT NULL PRIMARY KEY,
    org_id         TEXT        NOT NULL,
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down

DROP TABLE IF EXISTS cre.org_resolver_cache;
