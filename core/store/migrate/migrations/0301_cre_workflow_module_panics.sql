-- +goose Up

CREATE TABLE IF NOT EXISTS cre.workflow_module_panics (
    key         TEXT        NOT NULL PRIMARY KEY,
    val         BYTEA       NOT NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down

DROP TABLE IF EXISTS cre.workflow_module_panics;
