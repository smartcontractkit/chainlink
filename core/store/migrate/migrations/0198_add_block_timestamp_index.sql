-- +goose Up
create index log_poller_blocks_by_timestamp on evm.log_poller_blocks (evm_chain_id, block_timestamp);

-- +goose Down
DROP INDEX IF EXISTS evm.log_poller_blocks_by_timestamp;