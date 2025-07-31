-- +goose Up
-- +goose StatementBegin
-- unique constraint on workflow_owner and workflow_name
ALTER TABLE workflow_specs DROP CONSTRAINT unique_workflow_owner_name;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workflow_specs ADD CONSTRAINT unique_workflow_owner_name unique (workflow_owner, workflow_name);
-- +goose StatementEnd