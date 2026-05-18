-- +goose Up
CREATE TABLE logs (
  evm_chain_id numeric(78,0) NOT NULL REFERENCES evm_chains (id) DEFERRABLE,
  log_index bigint NOT NULL,
  block_hash bytea NOT NULL,
  block_number bigint NOT NULL CHECK (block_number > 0),
  address bytea NOT NULL,
  event_sig bytea NOT NULL,
  topics bytea[] NOT NULL,
  tx_hash bytea NOT NULL,
  data bytea NOT NULL,
  created_at timestamptz NOT NULL,
  PRIMARY KEY (block_hash, log_index, evm_chain_id)
);

-- Hot path query - clients searching for their logs.
CREATE INDEX logs_idx ON logs(evm_chain_id, block_number, address, event_sig);

CREATE TABLE log_poller_blocks (
    evm_chain_id numeric(78,0) NOT NULL REFERENCES evm_chains (id) DEFERRABLE,
    block_hash bytea NOT NULL,
    block_number bigint NOT NULL CHECK (block_number > 0),
    created_at timestamptz NOT NULL,
    -- Only permit one block_number at a time
    -- i.e. the poller is only ever aware of the canonical branch
    PRIMARY KEY (block_number, evm_chain_id)
);

-- +goose Down
DROP INDEX logs_idx;
DROP TABLE logs;
DROP TABLE log_poller_blocks;
