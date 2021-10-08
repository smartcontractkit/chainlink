-- +goose Up
ALTER TABLE pipeline_runs
    RENAME COLUMN errors TO fatal_errors;
ALTER TABLE pipeline_runs
    ADD COLUMN all_errors jsonb;

-- +goose Down
ALTER TABLE pipeline_runs
    RENAME COLUMN fatal_errors TO errors;
ALTER TABLE pipeline_runs
    DROP COLUMN all_errors;
