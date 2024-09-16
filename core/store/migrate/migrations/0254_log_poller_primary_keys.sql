-- +goose Up

ALTER TABLE evm.logs DROP CONSTRAINT logs_pkey;
ALTER TABLE evm.logs ADD COLUMN id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY;
CREATE UNIQUE INDEX idx_logs_chain_block_logindex ON evm.logs (evm_chain_id, block_number, log_index);
ALTER TABLE evm.log_poller_blocks DROP CONSTRAINT log_poller_blocks_pkey;
ALTER TABLE evm.log_poller_blocks ADD COLUMN id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY;
DROP INDEX IF EXISTS evm.idx_evm_log_poller_blocks_order_by_block;
DROP INDEX IF EXISTS evm.idx_evm_log_poller_blocks_block_number_evm_chain_id;
CREATE UNIQUE INDEX idx_log_poller_blocks_chain_block ON evm.log_poller_blocks (evm_chain_id, block_number DESC);
DROP INDEX IF EXISTS evm.idx_evm_logs_ordered_by_block_and_created_at;
CREATE INDEX idx_logs_chain_address_event_block_logindex ON evm.logs (evm_chain_id, address, event_sig, block_number, log_index);

-- +goose Down

DROP INDEX IF EXISTS evm.idx_logs_chain_address_event_block_logindex;
CREATE INDEX idx_evm_logs_ordered_by_block_and_created_at ON evm.logs (evm_chain_id, address, event_sig, block_number, created_at);
DROP INDEX IF EXISTS evm.idx_log_poller_blocks_chain_block;
ALTER TABLE evm.log_poller_blocks DROP COLUMN id;
ALTER TABLE evm.log_poller_blocks ADD PRIMARY KEY (block_number, evm_chain_id);
DROP INDEX IF EXISTS evm.idx_logs_chain_block_logindex;
ALTER TABLE evm.logs DROP COLUMN id;
ALTER TABLE evm.logs ADD PRIMARY KEY (block_hash, log_index, evm_chain_id);

