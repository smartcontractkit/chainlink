-- +goose Up
ALTER TABLE evm.txes ADD COLUMN "signal_callback" BOOL DEFAULT FALSE;
ALTER TABLE evm.txes ADD COLUMN "callback_completed" BOOL DEFAULT FALSE;

UPDATE evm.txes
SET signal_callback = TRUE AND callback_completed = FALSE
WHERE evm.txes.pipeline_task_run_id IN (
    SELECT pipeline_task_runs.id FROM pipeline_task_runs
    INNER JOIN pipeline_runs ON pipeline_runs.id = pipeline_task_runs.pipeline_run_id
    WHERE pipeline_runs.state = 'suspended'
);

-- +goose Down
ALTER TABLE evm.txes DROP COLUMN "signal_callback";
ALTER TABLE evm.txes DROP COLUMN "callback_completed";
