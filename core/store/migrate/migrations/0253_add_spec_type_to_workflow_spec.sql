-- +goose Up
-- +goose StatementBegin

ALTER TABLE workflow_specs ADD COLUMN spec_type varchar(255) DEFAULT 'yaml';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE workflow_specs DROP COLUMN spec_type;

-- +goose StatementEnd