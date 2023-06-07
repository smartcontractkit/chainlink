-- +goose Up
-- +goose StatementBegin
ALTER TABLE jobs
ADD COLUMN type_spec JSONB NOT NULL;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE jobs DROP COLUMN type_spec;
-- +goose StatementEnd
