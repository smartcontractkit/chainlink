-- +goose Up
-- +goose StatementBegin

-- The table was created as public.evm_upkeep_states by 0189 and moved to
-- evm.upkeep_states by 0194. Dropping the table also drops its indexes.
DROP TABLE IF EXISTS evm.upkeep_states;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

CREATE TABLE evm.upkeep_states (
  id SERIAL,
  work_id TEXT NOT NULL,
  evm_chain_id NUMERIC(20) NOT NULL,
  upkeep_id NUMERIC(78) NOT NULL, -- upkeep id is an evm word (uint256) which has a max size of precision 78
  completion_state SMALLINT NOT NULL,
  ineligibility_reason SMALLINT NOT NULL,
  block_number BIGINT NOT NULL,
  inserted_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
  -- named to match the constraint 0189 created implicitly, before 0194 renamed the table
  CONSTRAINT evm_upkeep_states_pkey PRIMARY KEY (id),
  CONSTRAINT work_id_len_chk CHECK (
    length(work_id) > 0 AND length(work_id) < 255
  )
);

CREATE UNIQUE INDEX idx_evm_upkeep_state_chainid_workid ON evm.upkeep_states (evm_chain_id, work_id);
CREATE INDEX idx_evm_upkeep_state_added_at_chain_id ON evm.upkeep_states (evm_chain_id, inserted_at);

-- +goose StatementEnd
