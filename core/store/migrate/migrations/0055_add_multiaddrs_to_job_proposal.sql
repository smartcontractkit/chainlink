-- +goose Up
ALTER TABLE job_proposals
ADD COLUMN multiaddrs TEXT[] DEFAULT NULL;
-- +goose Down
ALTER TABLE job_proposals
DROP COLUMN multiaddrs;
