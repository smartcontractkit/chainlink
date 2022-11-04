-- +goose Up
ALTER TABLE eth_txes DROP CONSTRAINT eth_txes_evm_chain_id_from_address_fkey, ADD CONSTRAINT eth_txes_evm_chain_id_from_address_fkey FOREIGN KEY (evm_chain_id, from_address) REFERENCES public.evm_key_states(evm_chain_id, address) ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE NOT VALID;
-- +goose Down
ALTER TABLE eth_txes DROP CONSTRAINT eth_txes_evm_chain_id_from_address_fkey, ADD CONSTRAINT eth_txes_evm_chain_id_from_address_fkey FOREIGN KEY (evm_chain_id, from_address) REFERENCES public.evm_key_states(evm_chain_id, address) DEFERRABLE INITIALLY IMMEDIATE NOT VALID;
