-- +goose Up
-- +goose StatementBegin

ALTER TABLE workflow_specs ADD COLUMN config varchar(255) DEFAULT '';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE workflow_specs DROP COLUMN config;

-- +goose StatementEnd