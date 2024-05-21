-- +goose Up
ALTER TABLE flux_monitor_specs ADD COLUMN drumbeat_random_delay bigint NOT NULL DEFAULT 0;

	UPDATE flux_monitor_specs SET drumbeat_schedule = '' where drumbeat_schedule IS NULL;
	ALTER TABLE flux_monitor_specs ALTER COLUMN drumbeat_schedule SET DEFAULT '';
	ALTER TABLE flux_monitor_specs ALTER COLUMN drumbeat_schedule SET NOT NULL;

-- +goose Down
ALTER TABLE flux_monitor_specs ALTER COLUMN drumbeat_schedule SET NULL;
ALTER TABLE flux_monitor_specs ALTER COLUMN drumbeat_schedule DROP DEFAULT;
ALTER TABLE flux_monitor_specs DROP COLUMN drumbeat_random_delay;
