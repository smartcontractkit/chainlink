-- +goose Up
ALTER TABLE bridge_types ADD COLUMN IF NOT EXISTS use_connection_manager boolean NOT NULL DEFAULT false;
-- +goose Down
ALTER TABLE bridge_types DROP COLUMN IF EXISTS use_connection_manager;
