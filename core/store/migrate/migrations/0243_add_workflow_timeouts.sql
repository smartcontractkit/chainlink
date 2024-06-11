-- +goose Up
-- +goose StatementBegin
ALTER TABLE workflow_specs ADD COLUMN max_execution_duration bigint;
ALTER TABLE workflow_specs ADD COLUMN max_step_duration bigint;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workflow_specs DROP COLUMN max_execution_duration;
ALTER TABLE workflow_specs DROP COLUMN max_step_duration;
-- +goose StatementEnd
