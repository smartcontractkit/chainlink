-- +goose Up
-- +goose StatementBegin

-- We must remove the old contraint to add an enum value to support Postgres v11
ALTER TABLE job_proposals
DROP CONSTRAINT chk_job_proposals_status_fsm;

-- Drop the cancelled enum value. Unfortunately postgres does not support a
-- a way to remove a value from an enum.
ALTER TYPE job_proposal_status RENAME TO job_proposal_status_old;
CREATE TYPE job_proposal_status AS ENUM('pending', 'approved', 'rejected', 'cancelled');

ALTER TABLE job_proposals ALTER COLUMN status TYPE job_proposal_status USING status::text::job_proposal_status;

DROP TYPE job_proposal_status_old;

-- Add the contraint back
ALTER TABLE job_proposals
ADD CONSTRAINT chk_job_proposals_status_fsm CHECK (
	(status = 'pending' AND external_job_id IS NULL) OR
	(status = 'approved' AND external_job_id IS NOT NULL) OR
	(status = 'rejected' AND external_job_id IS NULL)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- We must remove the old contraint to remove an enum value
ALTER TABLE job_proposals
DROP CONSTRAINT chk_job_proposals_status_fsm;

-- Drop the cancelled enum value. Unfortunately postgres does not support a
-- a way to remove a value from an enum.
ALTER TYPE job_proposal_status RENAME TO job_proposal_status_old;
CREATE TYPE job_proposal_status AS ENUM('pending', 'approved', 'rejected');

-- This will fail if any records are using the 'cancelled' enum.
-- Manually update these as we cannot decide what you want to do with them.
--
ALTER TABLE job_proposals ALTER COLUMN status TYPE job_proposal_status USING status::text::job_proposal_status;

DROP TYPE job_proposal_status_old;

-- Add the contraint back
ALTER TABLE job_proposals
ADD CONSTRAINT chk_job_proposals_status_fsm CHECK (
	(status = 'pending' AND external_job_id IS NULL) OR
	(status = 'approved' AND external_job_id IS NOT NULL) OR
	(status = 'rejected' AND external_job_id IS NULL)
);

-- +goose StatementEnd
