-- +goose Up
DELETE FROM node_versions WHERE version IN (
    SELECT version FROM node_versions ORDER BY created_at DESC OFFSET 1
);
CREATE UNIQUE INDEX idx_only_one_node_version ON node_versions ((version IS NOT NULL));

-- +goose Down
DROP INDEX idx_only_one_node_version;
