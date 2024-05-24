-- +goose Up
-- +goose StatementBegin
ALTER TABLE workflow_executions
DROP CONSTRAINT workflow_executions_workflow_id_fkey,
ADD CONSTRAINT workflow_executions_workflow_id_fkey
   FOREIGN KEY (workflow_id)
   REFERENCES workflow_specs(workflow_id)
   ON DELETE CASCADE;

ALTER TABLE workflow_steps
DROP CONSTRAINT workflow_steps_workflow_execution_id_fkey,
ADD CONSTRAINT workflow_steps_workflow_execution_id_fkey
   FOREIGN KEY (workflow_execution_id)
   REFERENCES workflow_executions(id)
   ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workflow_executions
DROP CONSTRAINT workflow_executions_workflow_id_fkey,
ADD CONSTRAINT workflow_executions_workflow_id_fkey
   FOREIGN KEY (workflow_id)
   REFERENCES workflow_specs(workflow_id);

ALTER TABLE workflow_steps
DROP CONSTRAINT workflow_steps_workflow_execution_id_fkey,
ADD CONSTRAINT workflow_steps_workflow_execution_id_fkey
   FOREIGN KEY (workflow_execution_id)
   REFERENCES workflow_executions(id);
-- +goose StatementEnd
