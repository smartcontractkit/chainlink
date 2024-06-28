INSERT INTO evm.log_poller_blocks (evm_chain_id, block_hash, block_number, created_at, block_timestamp, finalized_block_number)
SELECT '{{ .ChainID }}', block_hash, block_number, created_at, block_timestamp, finalized_block_number FROM {{ .Schema }}.log_poller_blocks;

DROP TABLE {{ .Schema }}.log_poller_blocks;
