-- +goose Up
ALTER TABLE job_proposals
ADD COLUMN remote_uuid UUID NOT NULL;

CREATE UNIQUE INDEX idx_job_proposals_remote_uuid ON job_proposals(remote_uuid);
-- +goose Down
ALTER TABLE job_proposals
DROP COLUMN remote_uuid;
