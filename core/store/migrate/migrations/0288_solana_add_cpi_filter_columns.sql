-- +goose Up
ALTER TABLE solana.log_poller_filters 
ADD COLUMN extra_filter_config JSONB;

-- +goose Down
ALTER TABLE solana.log_poller_filters 
DROP COLUMN IF EXISTS extra_filter_config;
