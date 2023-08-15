-- +goose Up

CREATE TABLE evm_upkeep_states (
    work_id TEXT PRIMARY KEY,
    evm_chain_id NUMERIC NOT NULL,
    upkeep_id BYTEA NOT NULL,
    completion_state SMALLINT NOT NULL,
    ineligibility_reason NUMERIC NOT NULL,
    block_number NUMERIC,
    added_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_evm_upkeep_state_chainid_workid ON evm_upkeep_states (evm_chain_id, work_id);
CREATE INDEX idx_evm_upkeep_state_added_at_chain_id ON evm_upkeep_states (evm_chain_id, added_at);

-- +goose Down

DROP INDEX IF EXISTS idx_evm_upkeep_state_chainid_workid;
DROP INDEX IF EXISTS idx_evm_upkeep_state_added_at_chain_id;

DROP TABLE evm_upkeep_states;