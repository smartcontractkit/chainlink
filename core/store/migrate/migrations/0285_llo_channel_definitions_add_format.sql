-- +goose Up
-- +goose StatementBegin
ALTER TABLE
  channel_definitions
ADD
  COLUMN format bigint;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
  channel_definitions DROP COLUMN format;
-- +goose StatementEnd
