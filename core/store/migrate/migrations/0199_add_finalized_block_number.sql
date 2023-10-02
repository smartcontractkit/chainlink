-- +goose Up
ALTER TABLE evm.log_poller_blocks
    ADD COLUMN last_finalized_block_number
        bigint not null
        default 0
        check (last_finalized_block_number >= 0);


-- +goose Down
ALTER TABLE evm.log_poller_blocks
    DROP COLUMN last_finalized_block_number;
