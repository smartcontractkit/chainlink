-- +goose Up
--- Remove all but most recently added contract_address for each chain. We will no longer allow duplicates, but enforcing that with a db constraint requires CREATE OPERATOR (admin) privilege
DELETE FROM ocr_oracle_specs WHERE id IN (SELECT id FROM (SELECT id, MAX(id) OVER(PARTITION BY evm_chain_id, contract_address ORDER BY id) AS max FROM ocr_oracle_specs) x WHERE id != max);

-- +goose Down
DROP INDEX IF EXISTS ocr_oracle_specs_unique_contract_addr;
DROP OPERATOR CLASS IF EXISTS wildcard_cmp USING BTREE CASCADE;
DROP FUNCTION IF EXISTS wildcard_cmp(INTEGER, INTEGER) CASCADE;
CREATE UNIQUE INDEX IF NOT EXISTS unique_contract_addr ON ocr_oracle_specs (contract_address) WHERE evm_chain_id IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS unique_contract_addr_per_chain ON ocr_oracle_specs (contract_address, evm_chain_id) WHERE evm_chain_id IS NOT NULL;
