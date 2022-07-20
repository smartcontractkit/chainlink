-- +goose Up
ALTER TABLE job_proposals
RENAME COLUMN job_id TO external_job_id;

ALTER INDEX idx_job_proposals_job_id RENAME TO idx_job_proposals_external_job_id;

-- +goose Down
ALTER TABLE job_proposals
RENAME COLUMN external_job_id TO job_id;

ALTER INDEX idx_job_proposals_external_job_id RENAME TO idx_job_proposals_job_id;
