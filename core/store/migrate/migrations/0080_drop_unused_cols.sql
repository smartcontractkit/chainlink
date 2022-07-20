-- +goose Up
ALTER TABLE offchainreporting_oracle_specs DROP COLUMN monitoring_endpoint;
ALTER TABLE jobs ADD COLUMN created_at timestamptz;

UPDATE jobs SET created_at=offchainreporting_oracle_specs.created_at FROM offchainreporting_oracle_specs WHERE jobs.offchainreporting_oracle_spec_id = offchainreporting_oracle_specs.id;
UPDATE jobs SET created_at=direct_request_specs.created_at FROM direct_request_specs WHERE jobs.direct_request_spec_id = direct_request_specs.id;
UPDATE jobs SET created_at=flux_monitor_specs.created_at FROM flux_monitor_specs WHERE jobs.flux_monitor_spec_id = flux_monitor_specs.id;
UPDATE jobs SET created_at=keeper_specs.created_at FROM keeper_specs WHERE jobs.keeper_spec_id = keeper_specs.id;
UPDATE jobs SET created_at=cron_specs.created_at FROM cron_specs WHERE jobs.cron_spec_id = cron_specs.id;
UPDATE jobs SET created_at=vrf_specs.created_at FROM vrf_specs WHERE jobs.vrf_spec_id = vrf_specs.id;
UPDATE jobs SET created_at=webhook_specs.created_at FROM webhook_specs WHERE jobs.webhook_spec_id = webhook_specs.id;

UPDATE jobs SET created_at = NOW() WHERE created_at IS NULL;
CREATE INDEX idx_jobs_created_at ON jobs USING BRIN (created_at);
ALTER TABLE jobs ALTER COLUMN created_at SET NOT NULL;

-- +goose Down
ALTER TABLE offchainreporting_oracle_specs ADD COLUMN monitoring_endpoint text;
ALTER TABLE jobs DROP COLUMN created_at;
