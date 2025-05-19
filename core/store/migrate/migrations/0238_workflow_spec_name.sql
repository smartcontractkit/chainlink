-- +goose Up
-- +goose StatementBegin
ALTER TABLE workflow_specs ADD COLUMN workflow_name varchar(255);

-- ensure that we can forward migrate to non-null name
UPDATE workflow_specs
SET
    workflow_name = workflow_id
WHERE
    workflow_name IS NULL;

ALTER TABLE workflow_specs ALTER COLUMN workflow_name SET NOT NULL;

-- unique constraint on workflow_owner and workflow_name
ALTER TABLE workflow_specs ADD CONSTRAINT unique_workflow_owner_name unique (workflow_owner, workflow_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workflow_specs DROP CONSTRAINT unique_workflow_owner_name;
ALTER TABLE workflow_specs DROP COLUMN workflow_name;
-- +goose StatementEnd