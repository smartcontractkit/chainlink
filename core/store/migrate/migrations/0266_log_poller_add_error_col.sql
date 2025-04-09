-- +goose Up

-- Add error column to solana.logs
ALTER TABLE solana.logs ADD COLUMN error TEXT;
ALTER TABLE solana.log_poller_filters ADD COLUMN include_reverted BOOLEAN NOT NULL DEFAULT false;

-- +goose Down

-- drop error column
ALTER TABLE solana.logs DROP COLUMN error;
ALTER TABLE solana.log_poller_filters DROP COLUMN include_reverted;
