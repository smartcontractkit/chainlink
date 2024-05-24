-- +goose Up
CREATE INDEX idx_unfinished_pipeline_task_runs ON pipeline_task_runs (pipeline_run_id) WHERE finished_at IS NULL;
-- +goose Down
DROP INDEX idx_unfinished_pipeline_task_runs;
