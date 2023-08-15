-- +goose Up

CREATE TABLE evm_upkeep_state (
    work_id TEXT PRIMARY KEY,
    evm_chain_id NUMERIC NOT NULL,
    upkeep_id NUMERIC NOT NULL,
    completion_state SMALLINT NOT NULL,
    ineligibility_reason NUMERIC NOT NULL,
    block_number NUMERIC,
    added_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_evm_upkeep_state_chainid_workid ON evm_upkeep_state (evm_chain_id, work_id);

-- +goose Down

DROP INDEX IF EXISTS idx_evm_upkeep_state_chainid_workid;

DROP TABLE evm_upkeep_state;