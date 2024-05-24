-- +goose Up
CREATE UNIQUE INDEX idx_pipeline_task_runs_unique_task_spec_id_per_run ON pipeline_task_runs (pipeline_task_spec_id, pipeline_run_id);

-- +goose Down
DROP INDEX IF EXISTS idx_pipeline_task_runs_unique_task_spec_id_per_run;
