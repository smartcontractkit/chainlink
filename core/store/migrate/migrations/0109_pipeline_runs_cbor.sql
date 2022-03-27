-- +goose Up
ALTER TABLE pipeline_task_runs DROP COLUMN "output";
ALTER TABLE pipeline_task_runs ADD COLUMN "output" bytea;

ALTER TABLE pipeline_runs DROP COLUMN "outputs";
ALTER TABLE pipeline_runs ADD COLUMN "outputs" bytea;

ALTER TABLE pipeline_runs DROP COLUMN "inputs";
ALTER TABLE pipeline_runs ADD COLUMN "inputs" bytea;

-- +goose Down
ALTER TABLE pipeline_task_runs DROP COLUMN "output";
ALTER TABLE pipeline_task_runs ADD COLUMN "output" jsonb;

ALTER TABLE pipeline_runs DROP COLUMN "outputs";
ALTER TABLE pipeline_runs ADD COLUMN "outputs" jsonb;

ALTER TABLE pipeline_runs DROP COLUMN "inputs";
ALTER TABLE pipeline_runs ADD COLUMN "inputs" jsonb;
