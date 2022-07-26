-- +goose Up
-- errors column was renamed to fatal_errors, see 0076_add_non_fatal_errors_to_runs.sql
-- but the constraint pipeline_runs_check was not updated
	ALTER TABLE pipeline_runs DROP CONSTRAINT pipeline_runs_check;
	ALTER TABLE pipeline_runs ADD CONSTRAINT pipeline_runs_check CHECK (
		((state IN ('completed')) AND (finished_at IS NOT NULL) AND (num_nulls(outputs) = 0))
			OR 
        ((state IN ('errored')) AND (finished_at IS NOT NULL) AND (num_nulls(fatal_errors, all_errors) = 0))
			OR 
		((state IN ('running', 'suspended')) AND num_nulls(finished_at, outputs, fatal_errors) = 3)
	);

-- +goose Down
-- we cannot make a precise rollback, due to a wrong column name (errors => fatal_errors)
-- therefore the rollback flow will fix it for pre-0132 state as well...
	ALTER TABLE pipeline_runs DROP CONSTRAINT pipeline_runs_check;
	ALTER TABLE pipeline_runs ADD CONSTRAINT pipeline_runs_check CHECK (
		((state IN ('completed', 'errored')) AND (finished_at IS NOT NULL) AND (num_nulls(outputs, fatal_errors) = 0))
			OR 
		((state IN ('running', 'suspended')) AND num_nulls(finished_at, outputs, fatal_errors) = 3)
	);
