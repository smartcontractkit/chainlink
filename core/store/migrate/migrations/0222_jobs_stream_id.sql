-- +goose Up
ALTER TABLE jobs ADD COLUMN stream_id BIGINT UNIQUE;

-- +goose Down
ALTER TABLE jobs DROP COLUMN stream_id;

