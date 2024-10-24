-- +goose Up
-- +goose StatementBegin
ALTER TABLE feeds_managers
ADD COLUMN disabled_at TIMESTAMP WITH TIME ZONE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE feeds_managers
DROP COLUMN IF EXISTS disabled_at;
-- +goose StatementEnd
