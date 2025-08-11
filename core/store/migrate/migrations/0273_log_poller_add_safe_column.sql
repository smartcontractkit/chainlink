-- +goose Up
-- Add safe_block_number column to evm.log_poller_blocks
-- This column is used to track the safe block number for each chain to be consumed from the log poller
ALTER TABLE evm.log_poller_blocks
    ADD COLUMN safe_block_number
        bigint not null
        default 0
        check (safe_block_number >= 0);


-- +goose Down
ALTER TABLE evm.log_poller_blocks
    DROP COLUMN safe_block_number;
