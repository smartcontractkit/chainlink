-- +goose Up
DROP INDEX IF EXISTS unique_contract_addr_per_chain;
DROP INDEX IF EXISTS unique_contract_addr;
CREATE FUNCTION wildcard_cmp(INTEGER, INTEGER) RETURNS INTEGER AS $$ SELECT CASE $1 = -1 OR $2 = -1 WHEN TRUE THEN 0 ELSE ($1-$2) END$$ PARALLEL SAFE IMMUTABLE LANGUAGE SQL;
CREATE OPERATOR CLASS wildcard_cmp FOR TYPE INTEGER USING BTREE AS FUNCTION 1 wildcard_cmp (INTEGER, INTEGER), OPERATOR 1 <, OPERATOR 2 <=, OPERATOR 3 =, OPERATOR 4 >=, OPERATOR 5 >;
-- Remove all but most recently added contract_address for each chain, before creating UNIQUE constraint
DELETE FROM ocr_oracle_specs WHERE id IN (SELECT id FROM (SELECT id, MAX(id) OVER(PARTITION BY evm_chain_id, contract_address ORDER BY id) AS max FROM ocr_oracle_specs) x WHERE id != max);
CREATE UNIQUE INDEX ocr_oracle_specs_unique_contract_addr ON ocr_oracle_specs (COALESCE(evm_chain_id::INTEGER, -1) wildcard_cmp, contract_address);

-- +goose Down
DROP INDEX IF EXISTS ocr_oracle_specs_unique_contract_addr;
DROP OPERATOR CLASS IF EXISTS wildcard_cmp USING BTREE CASCADE;
DROP FUNCTION IF EXISTS wildcard_cmp(INTEGER, INTEGER) CASCADE;
CREATE UNIQUE INDEX unique_contract_addr ON ocr_oracle_specs (contract_address) WHERE evm_chain_id IS NULL;
CREATE UNIQUE INDEX unique_contract_addr_per_chain ON ocr_oracle_specs (contract_address, evm_chain_id) WHERE evm_chain_id IS NOT NULL;
