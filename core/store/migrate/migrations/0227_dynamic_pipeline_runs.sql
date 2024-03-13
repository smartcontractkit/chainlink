-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.job_pipeline_specs (
    job_id INT NOT NULL,
    pipeline_spec_id INT NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE, -- Cannot use `primary` as column name in Postgres since it is a reserved keyword.
    CONSTRAINT pk_job_pipeline_spec PRIMARY KEY (job_id, pipeline_spec_id),
    CONSTRAINT fk_job FOREIGN KEY (job_id) REFERENCES public.jobs(id),
    CONSTRAINT fk_pipeline_spec FOREIGN KEY (pipeline_spec_id) REFERENCES public.pipeline_specs(id)
);

CREATE UNIQUE INDEX idx_unique_primary_per_job ON job_pipeline_spec(job_id) WHERE is_primary;

-- The moment this runs, we only have one job+pipeline_spec combination per job, complying with the unique index.
INSERT INTO public.job_pipeline_specs (job_id, pipeline_spec_id, is_primary)
SELECT id, pipeline_spec_id, TRUE
FROM public.jobs;


ALTER TABLE public.jobs DROP COLUMN pipeline_spec_id; -- Do we use CASCADE here? Does it have any relationship with other tables?
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.jobs ADD COLUMN pipeline_spec_id INT;

UPDATE public.jobs
SET pipeline_spec_id = jps.pipeline_spec_id
FROM public.job_pipeline_specs jps
WHERE jps.job_id = public.jobs.id
  AND jps.is_primary = TRUE;

DROP TABLE IF EXISTS public.job_pipeline_specs; -- Do we use CASCADE here?
-- +goose StatementEnd