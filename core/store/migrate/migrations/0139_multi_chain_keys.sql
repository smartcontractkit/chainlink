-- +goose Up
ALTER TABLE eth_key_states RENAME TO evm_key_states; -- Might as well rename it while we are here
CREATE UNIQUE INDEX idx_evm_key_states_evm_chain_id_address ON evm_key_states (evm_chain_id, address); -- it is now only unique per-chain
ALTER TABLE eth_txes DROP CONSTRAINT eth_txes_from_address_fkey; -- foreign key is now composite of chain id/address
ALTER TABLE evm_key_states DROP CONSTRAINT eth_key_states_address_key;
ALTER TABLE evm_key_states RENAME is_funding TO disabled; -- little hack here, we are removing is funding, to avoid accidentally sending from the wrong keys we disable the funding key
ALTER TABLE eth_txes ADD CONSTRAINT eth_txes_evm_chain_id_from_address_fkey FOREIGN KEY (evm_chain_id, from_address) REFERENCES evm_key_states (evm_chain_id, address) NOT VALID; -- not valid skips the check, this speeds things up and we know it's safe
CREATE INDEX idx_evm_key_states_address ON evm_key_states (address);

-- +goose Down
DROP INDEX idx_evm_key_states_address;
ALTER TABLE eth_txes DROP CONSTRAINT eth_txes_evm_chain_id_from_address_fkey;
ALTER TABLE evm_key_states RENAME disabled TO is_funding;
ALTER TABLE evm_key_states ADD CONSTRAINT eth_key_states_address_key UNIQUE (address);
ALTER TABLE eth_txes ADD CONSTRAINT eth_txes_from_address_fkey FOREIGN KEY (from_address) REFERENCES evm_key_states (address);
DROP INDEX idx_evm_key_states_evm_chain_id_address;
ALTER TABLE evm_key_states RENAME TO eth_key_states;

