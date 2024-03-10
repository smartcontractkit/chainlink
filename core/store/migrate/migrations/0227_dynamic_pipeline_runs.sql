-- +goose Up
-- +goose StatementBegin
CREATE TABLE job_pipeline_spec (
    job_id INT NOT NULL,
    pipeline_spec_id INT NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE, -- Cannot use `primary` as column name in Postgres since it is a reserved keyword.
    CONSTRAINT pk_job_pipeline_spec PRIMARY KEY (job_id, pipeline_spec_id),
    CONSTRAINT fk_job FOREIGN KEY (job_id) REFERENCES job(id),
    CONSTRAINT fk_pipeline_spec FOREIGN KEY (pipeline_spec_id) REFERENCES pipeline_spec(id)
);

INSERT INTO job_pipeline_spec (job_id, pipeline_spec_id, is_primary)
SELECT id, pipeline_spec_id, TRUE
FROM job;

ALTER TABLE job DROP COLUMN pipeline_spec_id; -- Do we use CASCADE here? Does it have any relationship with other tables?
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE job ADD COLUMN pipeline_spec_id INT;

UPDATE job
SET pipeline_spec_id = jps.pipeline_spec_id
FROM job_pipeline_spec jps
WHERE jps.job_id = job.id
  AND jps.is_primary = TRUE;

DROP TABLE IF EXISTS job_pipeline_spec; -- Do we use CASCADE here?
-- +goose StatementEnd