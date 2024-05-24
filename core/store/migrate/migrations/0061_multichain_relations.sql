-- +goose Up
ALTER TABLE evm_chains ADD COLUMN enabled BOOL DEFAULT TRUE NOT NULL;

ALTER TABLE eth_txes ADD COLUMN evm_chain_id numeric(78,0) REFERENCES evm_chains (id) DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE log_broadcasts ADD COLUMN evm_chain_id numeric(78,0) REFERENCES evm_chains (id) DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE heads ADD COLUMN evm_chain_id numeric(78,0) REFERENCES evm_chains (id) DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE eth_key_states ADD COLUMN evm_chain_id numeric(78,0) REFERENCES evm_chains (id) DEFERRABLE INITIALLY IMMEDIATE;

UPDATE eth_txes SET evm_chain_id = (SELECT id FROM evm_chains ORDER BY created_at, id ASC LIMIT 1);
UPDATE log_broadcasts SET evm_chain_id = (SELECT id FROM evm_chains ORDER BY created_at, id ASC LIMIT 1);
UPDATE heads SET evm_chain_id = (SELECT id FROM evm_chains ORDER BY created_at, id ASC LIMIT 1);
UPDATE eth_key_states SET evm_chain_id = (SELECT id FROM evm_chains ORDER BY created_at, id ASC LIMIT 1);

DROP INDEX IF EXISTS idx_eth_txes_min_unconfirmed_nonce_for_key;
DROP INDEX IF EXISTS idx_eth_txes_nonce_from_address;
DROP INDEX IF EXISTS idx_only_one_in_progress_tx_per_account;
DROP INDEX IF EXISTS idx_eth_txes_state_from_address;
DROP INDEX IF EXISTS idx_eth_txes_unstarted_subject_id;
CREATE INDEX idx_eth_txes_min_unconfirmed_nonce_for_key_evm_chain_id ON eth_txes(evm_chain_id, from_address, nonce) WHERE state = 'unconfirmed'::eth_txes_state;
CREATE UNIQUE INDEX idx_eth_txes_nonce_from_address_per_evm_chain_id ON eth_txes(evm_chain_id, from_address, nonce);
CREATE UNIQUE INDEX idx_only_one_in_progress_tx_per_account_id_per_evm_chain_id ON eth_txes(evm_chain_id, from_address) WHERE state = 'in_progress'::eth_txes_state;
CREATE INDEX idx_eth_txes_state_from_address_evm_chain_id ON eth_txes(evm_chain_id, from_address, state) WHERE state <> 'confirmed'::eth_txes_state;
CREATE INDEX idx_eth_txes_unstarted_subject_id_evm_chain_id ON eth_txes(evm_chain_id, subject, id) WHERE subject IS NOT NULL AND state = 'unstarted'::eth_txes_state;

DROP INDEX IF EXISTS idx_heads_hash;
DROP INDEX IF EXISTS idx_heads_number;
CREATE UNIQUE INDEX idx_heads_evm_chain_id_hash ON heads(evm_chain_id, hash);
CREATE INDEX idx_heads_evm_chain_id_number ON heads(evm_chain_id, number);

DROP INDEX IF EXISTS idx_log_broadcasts_unconsumed_job_id_v2;
DROP INDEX IF EXISTS log_consumptions_unique_v2_idx;
CREATE INDEX idx_log_broadcasts_unconsumed_job_id_v2 ON log_broadcasts(job_id, evm_chain_id) WHERE consumed = false AND job_id IS NOT NULL;
CREATE UNIQUE INDEX log_consumptions_unique_v2_idx ON log_broadcasts(job_id, block_hash, log_index, consumed, evm_chain_id) WHERE job_id IS NOT NULL;

ALTER TABLE eth_txes ALTER COLUMN evm_chain_id SET NOT NULL;
ALTER TABLE log_broadcasts ALTER COLUMN evm_chain_id SET NOT NULL;
ALTER TABLE heads ALTER COLUMN evm_chain_id SET NOT NULL;
ALTER TABLE eth_key_states ALTER COLUMN evm_chain_id SET NOT NULL;

ALTER TABLE vrf_specs ADD COLUMN evm_chain_id numeric(78,0) REFERENCES evm_chains (id) DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE direct_request_specs ADD COLUMN evm_chain_id numeric(78,0) REFERENCES evm_chains (id) DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE keeper_specs ADD COLUMN evm_chain_id numeric(78,0) REFERENCES evm_chains (id) DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE offchainreporting_oracle_specs ADD COLUMN evm_chain_id numeric(78,0) REFERENCES evm_chains (id) DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE flux_monitor_specs ADD COLUMN evm_chain_id numeric(78,0) REFERENCES evm_chains (id) DEFERRABLE INITIALLY IMMEDIATE;

-- +goose Down
ALTER TABLE evm_chains DROP COLUMN enabled;

ALTER TABLE heads DROP COLUMN evm_chain_id;
ALTER TABLE log_broadcasts DROP COLUMN evm_chain_id;
ALTER TABLE eth_txes DROP COLUMN evm_chain_id;
ALTER TABLE eth_key_states DROP COLUMN evm_chain_id;

CREATE UNIQUE INDEX idx_heads_hash ON heads(hash bytea_ops);
CREATE INDEX idx_heads_number ON heads(number int8_ops);

CREATE INDEX idx_eth_txes_min_unconfirmed_nonce_for_key ON eth_txes(from_address bytea_ops,nonce int8_ops) WHERE state = 'unconfirmed'::eth_txes_state;
CREATE UNIQUE INDEX idx_eth_txes_nonce_from_address ON eth_txes(from_address bytea_ops,nonce int8_ops);
CREATE UNIQUE INDEX idx_only_one_in_progress_tx_per_account ON eth_txes(from_address bytea_ops) WHERE state = 'in_progress'::eth_txes_state;
CREATE INDEX idx_eth_txes_state_from_address ON eth_txes(from_address bytea_ops,state enum_ops) WHERE state <> 'confirmed'::eth_txes_state;
CREATE INDEX idx_eth_txes_unstarted_subject_id ON eth_txes(subject uuid_ops,id int8_ops) WHERE subject IS NOT NULL AND state = 'unstarted'::eth_txes_state;

CREATE INDEX idx_log_broadcasts_unconsumed_job_id_v2 ON log_broadcasts(job_id int4_ops) WHERE consumed = false AND job_id IS NOT NULL;
CREATE UNIQUE INDEX log_consumptions_unique_v2_idx ON log_broadcasts(job_id int4_ops,block_hash bytea_ops,log_index int8_ops,consumed) WHERE job_id IS NOT NULL;

ALTER TABLE vrf_specs DROP COLUMN evm_chain_id;
ALTER TABLE direct_request_specs DROP COLUMN evm_chain_id;
ALTER TABLE keeper_specs DROP COLUMN evm_chain_id;
ALTER TABLE offchainreporting_oracle_specs DROP COLUMN evm_chain_id;
ALTER TABLE flux_monitor_specs DROP COLUMN evm_chain_id;
