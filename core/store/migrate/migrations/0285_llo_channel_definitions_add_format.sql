-- +goose Up
-- +goose StatementBegin
ALTER TABLE
  channel_definitions
ADD
  COLUMN format bigint NOT NULL DEFAULT 0;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
  channel_definitions DROP COLUMN format;
-- +goose StatementEnd
