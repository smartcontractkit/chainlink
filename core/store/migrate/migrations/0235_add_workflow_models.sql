-- +goose Up
-- +goose StatementBegin
CREATE TYPE workflow_status AS ENUM (
	'started',
	'errored',
	'timeout',
	'completed'
);

ALTER TABLE workflow_specs
       ADD CONSTRAINT fk_unique_workflow_id unique(workflow_id);

CREATE TABLE workflow_executions (
	id varchar(64) PRIMARY KEY,
	workflow_id varchar(64) references workflow_specs(workflow_id),
	status workflow_status NOT NULL,
	created_at timestamp with time zone,
	updated_at timestamp with time zone,
	finished_at timestamp with time zone
);

CREATE TABLE workflow_steps (
	id SERIAL PRIMARY KEY,
	workflow_execution_id varchar(64) references workflow_executions(id) NOT NULL,
	ref text NOT NULL,
	status workflow_status NOT NULL,
	inputs bytea,
	output_err text,
	output_value bytea,
	updated_at timestamp with time zone
);

ALTER TABLE workflow_steps
	ADD CONSTRAINT uniq_workflow_execution_id_ref unique(workflow_execution_id, ref);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workflow_steps
	DROP CONSTRAINT uniq_workflow_execution_id_ref;
DROP TABLE workflow_steps;
DROP TABLE workflow_executions;
ALTER TABLE workflow_specs
        DROP CONSTRAINT fk_unique_workflow_id;
DROP TYPE workflow_status;
-- +goose StatementEnd
