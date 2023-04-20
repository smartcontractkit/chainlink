-- +goose Up
-- +goose StatementBegin
ALTER TABLE evm_log_poller_filters ADD COLUMN retention BIGINT DEFAULT 0;
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
ALTER TABLE evm_log_poller_filters DROP COLUMN retention;
-- +goose StatementEnd
