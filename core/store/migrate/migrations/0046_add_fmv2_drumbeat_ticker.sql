-- +goose Up
ALTER TABLE flux_monitor_specs ADD COLUMN drumbeat_enabled boolean NOT NULL DEFAULT false;
ALTER TABLE flux_monitor_specs ADD COLUMN drumbeat_schedule text;
-- +goose Down
ALTER TABLE flux_monitor_specs DROP COLUMN drumbeat_enabled;
ALTER TABLE flux_monitor_specs DROP COLUMN drumbeat_schedule;
