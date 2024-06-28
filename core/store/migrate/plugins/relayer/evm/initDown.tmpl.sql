-- Do nothing for `evm` schema for backward compatibility
{{ if ne .Schema "evm"}}
/*
DROP TABLE {{ .Schema }}.receipts
DROP TABLE {{ .Schema }}.tx_attempts
DROP TABLE {{ .Schema }}.upkeep_states
DROP TABLE {{ .Schema }}.txes
DROP TABLE {{ .Schema }}.logs
DROP TABLE {{ .Schema }}.log_poller_filters
DROP TABLE {{ .Schema }}.log_poller_blocks
DROP TABLE {{ .Schema }}.key_states;
DROP TABLE {{ .Schema }}.heads;
*/

-- Copy data from old table to new table
INSERT INTO evm.forwarders (address, created_at, updated_at, evm_chain_id)
SELECT address, created_at, updated_at, '{{ .ChainID }}' 
FROM {{ .Schema }}.forwarders;

DROP TABLE {{ .Schema }}.forwarders;
{{ end}}