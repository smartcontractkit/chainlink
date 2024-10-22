-- +goose Up
-- +goose StatementBegin
ALTER TABLE feeds_managers
ADD COLUMN is_enabled BOOLEAN DEFAULT TRUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE feeds_managers
DROP COLUMN is_enabled;
-- +goose StatementEnd
