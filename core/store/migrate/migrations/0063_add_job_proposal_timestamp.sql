-- +goose Up
ALTER TABLE job_proposals
    ADD COLUMN proposed_at TIMESTAMP WITH TIME ZONE;

UPDATE job_proposals
    SET proposed_at = created_at;

ALTER TABLE job_proposals
    ALTER COLUMN proposed_at SET NOT NULL;

-- +goose Down
ALTER TABLE job_proposals DROP COLUMN proposed_at;
