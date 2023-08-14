-- +goose Up

CREATE TABLE evm_upkeep_state (
    id BIGSERIAL PRIMARY KEY,
    evm_chain_id NUMERIC NOT NULL,
    work_id TEXT NOT NULL,
    completion_state SMALLINT NOT NULL,
    block_number NUMERIC,
    added_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_evm_upkeep_state_chainid_workid ON evm_upkeep_state (evm_chain_id, work_id);

-- ALTER TABLE evm_upkeep_state 
-- ADD CONSTRAINT unique_chainid_workid 
-- UNIQUE USING INDEX idx_evm_upkeep_state_chainid_workid;

-- +goose Down

-- ALTER TABLE evm_upkeep_state DROP CONSTRAINT unique_chainid_workid;

DROP INDEX IF EXISTS idx_evm_upkeep_state_chainid_workid;

DROP TABLE evm_upkeep_state;