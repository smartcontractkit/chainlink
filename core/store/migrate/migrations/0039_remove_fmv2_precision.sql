-- +goose Up
ALTER TABLE flux_monitor_specs DROP COLUMN precision;
-- +goose Down
ALTER TABLE flux_monitor_specs ADD COLUMN precision integer;
