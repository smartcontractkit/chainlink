-- +goose Up

CREATE TABLE evm_upkeep_states (
  id SERIAL PRIMARY KEY,
  work_id TEXT NOT NULL,
  evm_chain_id NUMERIC(20) NOT NULL,
  upkeep_id NUMERIC(78) NOT NULL, -- upkeep id is an evm word (uint256) which has a max size of precision 78
  completion_state SMALLINT NOT NULL,
  ineligibility_reason SMALLINT NOT NULL,
  block_number BIGINT NOT NULL,
  inserted_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
  CONSTRAINT work_id_len_chk CHECK (
    length(work_id) > 0 AND length(work_id) < 255
  )
);

CREATE UNIQUE INDEX idx_evm_upkeep_state_chainid_workid ON evm_upkeep_states (evm_chain_id, work_id);
CREATE INDEX idx_evm_upkeep_state_added_at_chain_id ON evm_upkeep_states (evm_chain_id, inserted_at);

-- +goose Down

DROP INDEX IF EXISTS idx_evm_upkeep_state_chainid_workid;
DROP INDEX IF EXISTS idx_evm_upkeep_state_added_at_chain_id;

DROP TABLE evm_upkeep_states;