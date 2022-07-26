-- +goose Up
-- errors column was renamed to fatal_errors, see 0076_add_non_fatal_errors_to_runs.sql
-- but the constraint pipeline_runs_check was not updated
	ALTER TABLE pipeline_runs DROP CONSTRAINT pipeline_runs_check;
	ALTER TABLE pipeline_runs ADD CONSTRAINT pipeline_runs_check CHECK (
		((state IN ('completed')) AND (finished_at IS NOT NULL) AND (num_nulls(outputs) = 0) AND (num_nulls(fatal_errors) = 1))
			OR 
        ((state IN ('errored')) AND (finished_at IS NOT NULL) AND (num_nulls(fatal_errors, all_errors) = 0))
			OR 
		((state IN ('running', 'suspended')) AND num_nulls(finished_at, outputs, fatal_errors) = 3)
	);

-- +goose Down
	ALTER TABLE pipeline_runs DROP CONSTRAINT pipeline_runs_check;
	ALTER TABLE pipeline_runs ADD CONSTRAINT pipeline_runs_check CHECK (
		((state IN ('completed', 'errored')) AND (finished_at IS NOT NULL) AND (num_nulls(outputs, errors) = 0))
			OR 
		((state IN ('running', 'suspended')) AND num_nulls(finished_at, outputs, errors) = 3)
	);
