-- +goose Up

-- Start with dropping the index introduced in a previous migration - we are not going to use it
DROP INDEX IF EXISTS evm.log_poller_blocks_by_timestamp;

CREATE INDEX evm_logs_by_timestamp
    ON evm.logs (evm_chain_id, address, event_sig, block_timestamp, block_number);

-- +goose Down
create index log_poller_blocks_by_timestamp on evm.log_poller_blocks (evm_chain_id, block_timestamp);

drop index if exists evm.evm_logs_by_timestamp;



