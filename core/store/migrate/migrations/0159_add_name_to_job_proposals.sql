-- +goose Up
-- +goose StatementBegin

-- Add the name column to job proposals
ALTER TABLE job_proposals
ADD COLUMN name TEXT;

-- Attempt to populate the name field from a proposal's job spec definition.
-- If it does not match the regex, it will continue to search through the
-- versions to find one that matches. If none match, the job proposal name is
-- left blank.
UPDATE job_proposals
SET name = specs.name
FROM (
	SELECT job_proposal_id, (regexp_matches(job_proposal_specs.definition, 'name = ''(.+?)\''\n'))[1] as name, MAX(version)
	FROM job_proposal_specs
	GROUP BY job_proposal_id, name
) AS specs
WHERE job_proposals.id = specs.job_proposal_id

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE job_proposals
DROP COLUMN name;

-- +goose StatementEnd
