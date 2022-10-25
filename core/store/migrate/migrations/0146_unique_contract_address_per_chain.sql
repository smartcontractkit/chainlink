-- +goose Up

-- no-op

-- +goose Down
DROP INDEX IF EXISTS ocr_oracle_specs_unique_contract_addr;
DROP OPERATOR CLASS IF EXISTS wildcard_cmp USING BTREE CASCADE;
DROP FUNCTION IF EXISTS wildcard_cmp(INTEGER, INTEGER) CASCADE;
CREATE UNIQUE INDEX IF NOT EXISTS unique_contract_addr ON ocr_oracle_specs (contract_address) WHERE evm_chain_id IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS unique_contract_addr_per_chain ON ocr_oracle_specs (contract_address, evm_chain_id) WHERE evm_chain_id IS NOT NULL;
