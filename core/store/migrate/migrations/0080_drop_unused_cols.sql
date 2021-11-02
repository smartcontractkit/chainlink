-- +goose Up
ALTER TABLE offchainreporting_oracle_specs DROP COLUMN monitoring_endpoint;
ALTER TABLE jobs ADD COLUMN created_at timestamptz;
UPDATE jobs SET created_at = '1970-01-01';
CREATE INDEX idx_jobs_created_at ON jobs USING BRIN (created_at);
ALTER TABLE jobs ALTER COLUMN created_at SET NOT NULL;

-- +goose Down
ALTER TABLE offchainreporting_oracle_specs ADD COLUMN monitoring_endpoint text;
ALTER TABLE jobs DROP COLUMN created_at;
