-- +goose Up
CREATE TABLE cron_specs (
    id SERIAL PRIMARY KEY,
    cron_schedule text NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);
ALTER TABLE jobs ADD COLUMN cron_spec_id INT REFERENCES cron_specs(id),
DROP CONSTRAINT chk_only_one_spec,
ADD CONSTRAINT chk_only_one_spec CHECK (
    num_nonnulls(offchainreporting_oracle_spec_id, direct_request_spec_id, flux_monitor_spec_id, keeper_spec_id, cron_spec_id) = 1
);

-- +goose Down
ALTER TABLE jobs DROP CONSTRAINT chk_only_one_spec,
ADD CONSTRAINT chk_only_one_spec CHECK (
    num_nonnulls(offchainreporting_oracle_spec_id, direct_request_spec_id, flux_monitor_spec_id, keeper_spec_id) = 1
);

ALTER TABLE jobs DROP COLUMN cron_spec_id;
DROP TABLE IF EXISTS cron_specs;
