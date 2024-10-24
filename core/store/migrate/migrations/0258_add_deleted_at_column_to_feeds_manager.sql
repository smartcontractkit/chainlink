-- +goose Up
-- +goose StatementBegin
ALTER TABLE feeds_managers
ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE feeds_managers
DROP COLUMN IF EXISTS deleted_at;
-- +goose StatementEnd
