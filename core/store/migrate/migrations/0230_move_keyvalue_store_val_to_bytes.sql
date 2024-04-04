-- +goose Up

-- Add a new bytea column
ALTER TABLE job_kv_store ADD COLUMN val_bytea bytea;

-- Copy and convert the data from the jsonb column to the new bytea column
UPDATE job_kv_store SET val_bytea = convert_to(val::text, 'UTF8');

-- Drop the jsonb column
ALTER TABLE job_kv_store DROP COLUMN val;

-- +goose Down
ALTER TABLE job_kv_store ADD COLUMN val jsonb;
-- Bytea data may not be convertable to jsonb, so just drop the column
ALTER TABLE job_kv_store DROP COLUMN val_bytea;

