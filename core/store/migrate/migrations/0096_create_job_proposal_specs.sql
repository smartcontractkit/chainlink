-- +goose Up
-- +goose StatementBegin

-- Create a new enum type for the spec's status.
CREATE TYPE job_proposal_spec_status AS ENUM('pending', 'approved', 'rejected', 'cancelled');

-- Create a new table to store the versioned specs
CREATE TABLE job_proposal_specs (
    id SERIAL PRIMARY KEY,
    definition TEXT NOT NULL,
    version INTEGER NOT NULL,
    status job_proposal_spec_status NOT NULL,
    job_proposal_id INTEGER REFERENCES job_proposals(id),
    status_updated_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);

CREATE UNIQUE INDEX idx_job_proposals_job_proposal_id_and_version ON job_proposal_specs(job_proposal_id, version);
CREATE UNIQUE INDEX idx_job_proposal_specs_job_proposal_id_and_status ON job_proposal_specs(job_proposal_id) WHERE status = 'approved';

-- Seed existing data from job_proposals into the job proposal specs
INSERT INTO job_proposal_specs
(
    definition,
    version,
    status,
    job_proposal_id,
    status_updated_at,
    created_at,
    updated_at
)
SELECT spec,
       1,
       -- Cast to a string before casting to a job_proposal_spec_status because
       -- you can't cast from enum to enum. This is safe because the enums
       -- match exactly.
       status::varchar::job_proposal_spec_status,
       id,
       updated_at,
       proposed_at,
       updated_at
from job_proposals;

-- Update job proposals table with new fields
--   * Drop columns now that we have moved the data
--   * Add a pending update column
ALTER TABLE job_proposals
DROP COLUMN spec,
DROP COLUMN proposed_at,
ADD COLUMN pending_update BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Add the columns back into the job proposals table and drop the pending update
ALTER TABLE job_proposals
ADD COLUMN spec TEXT,
ADD COLUMN proposed_at timestamp with time zone,
DROP COLUMN pending_update;

-- Return the data back to the job_proposals table
UPDATE job_proposals
SET spec=jps.definition,
    proposed_at=jps.created_at,
    status=jps.status
FROM (
  SELECT a.definition, a.job_proposal_id, a.created_at, a.status::varchar::job_proposal_status
  FROM job_proposal_specs a
  INNER JOIN (
    SELECT job_proposal_id, MAX(version) ver
    FROM job_proposal_specs
    GROUP BY job_proposal_id
  ) b ON a.job_proposal_id = b.job_proposal_id AND a.version = b.ver
) AS jps
WHERE job_proposals.id = jps.job_proposal_id;

-- Add constraints to the new fields
ALTER TABLE job_proposals
ALTER COLUMN spec SET NOT NULL,
ALTER COLUMN proposed_at SET NOT NULL;

-- Drop the job_proposals table
DROP TABLE job_proposal_specs;

-- Drop the enum
DROP TYPE job_proposal_spec_status

-- +goose StatementEnd
