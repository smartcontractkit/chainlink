-- +goose Up
ALTER TABLE solana.log_poller_filters ADD COLUMN contract_idl TEXT;

-- +goose Down
ALTER TABLE solana.log_poller_filters DROP COLUMN contract_idl;

