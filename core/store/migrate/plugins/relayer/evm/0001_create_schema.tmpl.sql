-- +goose Up
CREATE SCHEMA IF NOT EXISTS {{.Schema}};
-- +goose Down
DROP SCHEMA IF EXISTS {{.Schema}};