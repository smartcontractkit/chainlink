-- +goose Up
ALTER TYPE workflow_status ADD VALUE 'completed_early_exit';

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
