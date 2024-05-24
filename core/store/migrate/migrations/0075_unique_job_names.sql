-- +goose Up

-- Rename duplicate jobs first
UPDATE jobs
SET name = jobs.name || ' (' || j.rank::text || ')'
FROM (SELECT id, row_number() OVER (PARTITION BY name ORDER BY id) AS rank FROM jobs) j
WHERE jobs.id = j.id AND j.rank > 1;

CREATE UNIQUE INDEX idx_jobs_name ON jobs (name);

-- +goose Down
DROP INDEX IF EXISTS idx_jobs_name;