-- +goose Up
-- +goose StatementBegin

-- Remove the unique contraint on the index
DROP INDEX idx_job_proposals_remote_uuid;
CREATE INDEX idx_job_proposals_remote_uuid
ON job_proposals(remote_uuid);

-- Create a unique partial index on approved job proposals
CREATE UNIQUE INDEX idx_approved_job_proposals_remote_uuid ON job_proposals(remote_uuid) WHERE (status = 'approved');

-- Create a unique partial index on pending job proposals
CREATE UNIQUE INDEX idx_pending_job_proposals_remote_uuid ON job_proposals(remote_uuid) WHERE (status = 'pending');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_pending_job_proposals_remote_uuid;
DROP INDEX idx_approved_job_proposals_remote_uuid;
DROP INDEX idx_job_proposals_remote_uuid;
CREATE UNIQUE INDEX idx_job_proposals_remote_uuid
ON job_proposals(remote_uuid);
-- +goose StatementEnd
