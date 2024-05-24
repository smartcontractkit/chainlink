-- +goose Up
ALTER TABLE pipeline_task_runs DROP CONSTRAINT chk_pipeline_task_run_fsm;
DELETE FROM pipeline_task_runs WHERE type = 'result';
ALTER TABLE pipeline_task_runs 
	ADD CONSTRAINT chk_pipeline_task_run_fsm CHECK (
		((finished_at IS NOT NULL) AND (num_nonnulls(output, error) != 2))
			OR 
		(num_nulls(finished_at, output, error) = 3)
	);

-- +goose Down
ALTER TABLE pipeline_task_runs DROP CONSTRAINT chk_pipeline_task_run_fsm;
ALTER TABLE pipeline_task_runs 
	ADD CONSTRAINT chk_pipeline_task_run_fsm CHECK (
		(((type <> 'result'::text) AND (((finished_at IS NULL) AND (error IS NULL) AND (output IS NULL)) 
			OR 
		((finished_at IS NOT NULL) AND (NOT ((error IS NOT NULL) AND (output IS NOT NULL)))))) 
			OR 
		((type = 'result'::text) AND (((output IS NULL) AND (error IS NULL) AND (finished_at IS NULL)) 
			OR 
		((output IS NOT NULL) AND (error IS NOT NULL) AND (finished_at IS NOT NULL))))));
