-- +goose Up
-- +goose StatementBegin
ALTER TABLE workflow_specs_v2 ADD COLUMN workflow_tag varchar(255) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workflow_specs_v2 DROP COLUMN workflow_tag;
-- +goose StatementEnd
