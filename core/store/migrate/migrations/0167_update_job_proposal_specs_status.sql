-- +goose Up
-- +goose StatementBegin
DROP INDEX idx_job_proposal_specs_job_proposal_id_and_status;

ALTER TYPE job_proposal_spec_status
RENAME TO job_proposal_spec_status_old;

CREATE TYPE job_proposal_spec_status AS ENUM(
    'pending',
    'approved',
    'rejected',
    'cancelled',
    'revoked'
);

ALTER TABLE job_proposal_specs
ALTER COLUMN status TYPE job_proposal_spec_status USING status::TEXT::job_proposal_spec_status;

DROP TYPE job_proposal_spec_status_old;

CREATE UNIQUE INDEX idx_job_proposal_specs_job_proposal_id_and_status ON job_proposal_specs(job_proposal_id)
WHERE status = 'approved';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_job_proposal_specs_job_proposal_id_and_status;

ALTER TYPE job_proposal_spec_status
RENAME TO job_proposal_spec_status_old;

CREATE TYPE job_proposal_spec_status AS ENUM('pending', 'approved', 'rejected', 'cancelled');

-- This will fail if any records are using the 'revoked' enum.
-- Manually update these as we cannot decide what you want to do with them.
ALTER TABLE job_proposal_specs
ALTER COLUMN status TYPE job_proposal_spec_status USING status::TEXT::job_proposal_spec_status;

DROP TYPE job_proposal_spec_status_old;

CREATE UNIQUE INDEX idx_job_proposal_specs_job_proposal_id_and_status ON job_proposal_specs(job_proposal_id)
WHERE status = 'approved';

-- +goose StatementEnd
