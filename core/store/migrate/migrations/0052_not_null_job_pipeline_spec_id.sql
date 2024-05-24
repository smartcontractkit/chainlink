-- +goose Up
ALTER TABLE jobs ALTER COLUMN pipeline_spec_id SET NOT NULL;
-- +goose Down
ALTER TABLE jobs ALTER COLUMN pipeline_spec_id DEFAULT NULL;
