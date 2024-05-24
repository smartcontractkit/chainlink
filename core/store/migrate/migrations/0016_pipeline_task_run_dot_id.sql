-- +goose Up
ALTER TABLE pipeline_task_runs ADD COLUMN dot_id text; 
UPDATE pipeline_task_runs SET dot_id = ts.dot_id FROM pipeline_task_specs ts WHERE ts.id = pipeline_task_runs.pipeline_task_spec_id;
ALTER TABLE pipeline_task_runs ALTER COLUMN dot_id SET NOT NULL, DROP COLUMN pipeline_task_spec_id;
DROP TABLE pipeline_task_specs;

CREATE UNIQUE INDEX ON pipeline_task_runs(pipeline_run_id, dot_id);
DROP INDEX idx_pipeline_task_runs_optimise_find_results;

-- +goose Down
ALTER TABLE pipeline_task_runs DROP COLUMN dot_id;
CREATE TABLE public.pipeline_task_specs (
    id BIGSERIAL PRIMARY KEY,
    dot_id text NOT NULL,
    pipeline_spec_id integer NOT NULL,
    type text NOT NULL,
    json jsonb NOT NULL,
    index integer DEFAULT 0 NOT NULL,
    successor_id integer,
    created_at timestamp with time zone NOT NULL
);
CREATE INDEX idx_pipeline_task_specs_created_at ON public.pipeline_task_specs USING brin (created_at);
CREATE INDEX idx_pipeline_task_specs_pipeline_spec_id ON public.pipeline_task_specs USING btree (pipeline_spec_id);
CREATE UNIQUE INDEX idx_pipeline_task_specs_single_output ON public.pipeline_task_specs USING btree (pipeline_spec_id) WHERE (successor_id IS NULL);
CREATE INDEX idx_pipeline_task_specs_successor_id ON public.pipeline_task_specs USING btree (successor_id);
CREATE INDEX idx_pipeline_task_runs_optimise_find_results ON public.pipeline_task_runs USING btree (pipeline_run_id);
