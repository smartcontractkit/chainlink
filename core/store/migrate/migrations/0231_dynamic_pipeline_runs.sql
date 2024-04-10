-- +goose Up
-- +goose StatementBegin
CREATE TABLE job_pipeline_specs (
    job_id INT NOT NULL,
    pipeline_spec_id INT NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT pk_job_pipeline_spec PRIMARY KEY (job_id, pipeline_spec_id),
    CONSTRAINT fk_job_pipeline_spec_job FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE DEFERRABLE,
    CONSTRAINT fk_job_pipeline_spec_pipeline_spec FOREIGN KEY (pipeline_spec_id) REFERENCES pipeline_specs(id) ON DELETE CASCADE DEFERRABLE
);

CREATE UNIQUE INDEX idx_unique_job_pipeline_spec_primary_per_job ON job_pipeline_specs(job_id) WHERE is_primary;

-- The moment this runs, we only have one job+pipeline_spec combination per job, complying with the unique index.
INSERT INTO job_pipeline_specs (job_id, pipeline_spec_id, is_primary)
SELECT id, pipeline_spec_id, TRUE
FROM jobs;

ALTER TABLE jobs DROP COLUMN pipeline_spec_id;

ALTER TABLE pipeline_runs ADD COLUMN pruning_key INT;

UPDATE pipeline_runs
SET pruning_key = pjps.job_id
FROM job_pipeline_specs pjps
WHERE pjps.pipeline_spec_id = pipeline_runs.pipeline_spec_id;

ALTER TABLE pipeline_runs ALTER COLUMN pruning_key SET NOT NULL;

ALTER TABLE pipeline_runs ADD CONSTRAINT fk_pipeline_runs_pruning_key FOREIGN KEY (pruning_key) REFERENCES jobs(id) ON DELETE CASCADE DEFERRABLE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE jobs ADD COLUMN pipeline_spec_id INT;

UPDATE jobs
SET pipeline_spec_id = jps.pipeline_spec_id
FROM job_pipeline_specs jps
WHERE jps.job_id = jobs.id
  AND jps.is_primary = TRUE;

ALTER TABLE pipeline_runs DROP COLUMN pruning_key;

DROP INDEX IF EXISTS idx_unique_primary_per_job;

DROP TABLE IF EXISTS job_pipeline_specs;
-- +goose StatementEnd