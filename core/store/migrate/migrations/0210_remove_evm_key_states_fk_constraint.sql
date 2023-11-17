-- +goose Up
ALTER TABLE evm.txes DROP CONSTRAINT eth_txes_evm_chain_id_from_address_fkey;
-- +goose Down
ALTER TABLE evm.txes ADD CONSTRAINT eth_txes_evm_chain_id_from_address_fkey FOREIGN KEY (evm_chain_id, from_address) REFERENCES evm.key_states(evm_chain_id, address) ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE NOT VALID;