-- +goose Up
CREATE INDEX logs_idx_block_number ON logs using brin(block_number);
CREATE INDEX logs_idx_evm_id_event_address_block ON logs using btree (evm_chain_id,event_sig,address,block_number);

-- +goose Down
DROP INDEX IF EXISTS logs_idx_block_number;
DROP INDEX IF EXISTS logs_idx_evm_id_event_address_block;
