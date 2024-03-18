-- +goose Up
ALTER TABLE jobs ADD COLUMN stream_id BIGINT;
CREATE UNIQUE INDEX idx_jobs_unique_stream_id ON jobs(stream_id) WHERE stream_id IS NOT NULL;

-- +goose Down
ALTER TABLE jobs DROP COLUMN stream_id;

