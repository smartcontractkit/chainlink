-- +goose Up
ALTER TABLE bridge_types ADD COLUMN use_connection_manager boolean NOT NULL DEFAULT false;
-- +goose Down
ALTER TABLE bridge_types DROP COLUMN use_connection_manager;
