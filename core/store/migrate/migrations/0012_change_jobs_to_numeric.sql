-- +goose Up
ALTER TABLE job_specs ALTER COLUMN min_payment TYPE numeric(78, 0) USING min_payment::numeric;
ALTER TABLE flux_monitor_specs ALTER COLUMN min_payment TYPE numeric(78, 0) USING min_payment::numeric;

-- +goose Down
ALTER TABLE job_specs ALTER COLUMN min_payment TYPE varchar(255) USING min_payment::varchar;
ALTER TABLE flux_monitor_specs ALTER COLUMN min_payment TYPE varchar(255) USING min_payment::varchar;
