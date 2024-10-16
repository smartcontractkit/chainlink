-- +goose Up
-- +goose StatementBegin
ALTER TABLE standardcapabilities_specs
ADD COLUMN oracle_factory JSONB;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE standardcapabilities_specs DROP COLUMN oracle_factory;
-- +goose StatementEnd
