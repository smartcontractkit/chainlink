-- +goose Up
CREATE SCHEMA IF NOT EXISTS {{.Schema}};
-- +goose Down
-- we don't know if the schema was created by this migration, so we can't drop it
