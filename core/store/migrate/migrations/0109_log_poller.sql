-- +goose Up
CREATE TABLE logs (
  evm_chain_id numeric(78,0) NOT NULL REFERENCES evm_chains (id) DEFERRABLE,
  log_index bigint NOT NULL,
  block_hash bytea NOT NULL,
  block_number bigint NOT NULL,
  address bytea NOT NULL,
  event_signature text NOT NULL,
  tx_hash bytea NOT NULL,
  data bytea NOT NULL,
  created_at timestamptz NOT NULL,
  PRIMARY KEY (block_hash, log_index, evm_chain_id)
);

CREATE TABLE log_poller_blocks (
    evm_chain_id numeric(78,0) NOT NULL REFERENCES evm_chains (id) DEFERRABLE,
    block_hash bytea NOT NULL,
    -- Only permit one block_number at a time
    -- i.e. the poller is only ever aware of the canonical branch
    block_number bigint UNIQUE NOT NULL,
    created_at timestamptz NOT NULL,
    PRIMARY KEY (block_hash, evm_chain_id)
);

-- +goose Down
DROP TABLE logs;
DROP TABLE log_poller_blocks;
