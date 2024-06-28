INSERT INTO evm.heads (hash, "number", parent_hash, created_at, "timestamp", l1_block_number, evm_chain_id, base_fee_per_gas) 
SELECT hash, "number", parent_hash, created_at, "timestamp", l1_block_number, '{{ .ChainID }}', base_fee_per_gas 
FROM {{ .Schema }}.heads;

DROP TABLE {{ .Schema }}.heads;
