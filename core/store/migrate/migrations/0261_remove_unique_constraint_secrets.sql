-- +goose Up
-- +goose StatementBegin
-- unique constraint on workflow_owner and workflow_name
ALTER TABLE workflow_specs DROP CONSTRAINT workflow_specs_secrets_id_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workflow_specs ADD CONSTRAINT workflow_specs_secrets_id_key unique (secrets_id);
-- +goose StatementEnd
