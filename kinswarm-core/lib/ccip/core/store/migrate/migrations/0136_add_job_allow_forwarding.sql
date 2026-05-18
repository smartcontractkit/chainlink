-- +goose Up
ALTER TABLE jobs ADD COLUMN allow_forwarding boolean DEFAULT FALSE;
-- +goose Down
ALTER TABLE jobs DROP COLUMN allow_forwarding;
