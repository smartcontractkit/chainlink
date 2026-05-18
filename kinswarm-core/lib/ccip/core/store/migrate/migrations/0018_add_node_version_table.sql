-- +goose Up
CREATE TABLE IF NOT EXISTS "node_versions" (
    "version" TEXT PRIMARY KEY,
    "created_at" timestamp without time zone NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS "node_versions";
