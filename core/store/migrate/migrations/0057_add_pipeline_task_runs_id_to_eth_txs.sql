-- +goose Up
ALTER TABLE eth_txes ADD COLUMN pipeline_task_run_id uuid UNIQUE;
ALTER TABLE eth_txes ADD COLUMN min_confirmations integer;
CREATE INDEX pipeline_runs_suspended ON pipeline_runs (id) WHERE state = 'suspended' ;

-- +goose Down
ALTER TABLE eth_txes DROP COLUMN pipeline_task_run_id;
ALTER TABLE eth_txes DROP COLUMN min_confirmations;
DROP INDEX pipeline_runs_suspended;
