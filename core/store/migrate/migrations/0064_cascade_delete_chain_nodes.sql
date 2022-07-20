-- +goose Up
ALTER TABLE nodes 
DROP CONSTRAINT nodes_evm_chain_id_fkey;

ALTER TABLE nodes
ADD FOREIGN KEY (evm_chain_id) REFERENCES evm_chains(id) ON DELETE CASCADE;

--+goose Down
ALTER TABLE nodes 
DROP CONSTRAINT nodes_evm_chain_id_fkey;

ALTER TABLE nodes
ADD FOREIGN KEY (evm_chain_id) REFERENCES evm_chains(id);
