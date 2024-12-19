-- +goose Up
-- +goose StatementBegin
-- Create the workflow_artifacts table
CREATE TABLE workflow_secrets (
    id SERIAL PRIMARY KEY,
    secrets_url TEXT,
    secrets_url_hash TEXT UNIQUE,
    contents TEXT
);

-- Create an index on the secrets_url_hash column
CREATE INDEX idx_secrets_url ON workflow_secrets(secrets_url);

-- Alter the workflow_specs table
ALTER TABLE workflow_specs
ADD COLUMN binary_url TEXT DEFAULT '',
ADD COLUMN config_url TEXT DEFAULT '',
ADD COLUMN secrets_id INT UNIQUE REFERENCES workflow_secrets(id) ON DELETE CASCADE;

-- Alter the config column type
ALTER TABLE workflow_specs
ALTER COLUMN config TYPE TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workflow_specs
DROP COLUMN IF EXISTS secrets_id,
DROP COLUMN IF EXISTS config_url,
DROP COLUMN IF EXISTS binary_url;

-- Change the config column back to character varying(255)
ALTER TABLE workflow_specs
ALTER COLUMN config TYPE CHARACTER VARYING(255);

-- Drop the index on the secrets_url_hash column
DROP INDEX IF EXISTS idx_secrets_url_hash;

-- Drop the workflow_artifacts table
DROP TABLE IF EXISTS workflow_secrets;
-- +goose StatementEnd

