-- +goose Up

ALTER TABLE pipeline_runs ALTER COLUMN meta DROP NOT NULL;


-- +goose Down

ALTER TABLE pipeline_runs ALTER COLUMN meta SET NOT NULL;

