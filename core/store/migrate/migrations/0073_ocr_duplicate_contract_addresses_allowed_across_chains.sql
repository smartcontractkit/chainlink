-- +goose Up

ALTER TABLE offchainreporting_oracle_specs DROP CONSTRAINT unique_contract_addr;
CREATE UNIQUE INDEX unique_contract_addr_per_chain ON offchainreporting_oracle_specs (contract_address, evm_chain_id) WHERE evm_chain_id IS NOT NULL;
CREATE UNIQUE INDEX unique_contract_addr ON offchainreporting_oracle_specs (contract_address) WHERE evm_chain_id IS NULL;

-- +goose Down

DROP INDEX unique_contract_addr;
DROP INDEX unique_contract_addr_per_chain;
ALTER TABLE offchainreporting_oracle_specs ADD CONSTRAINT unique_contract_addr UNIQUE (contract_address);
