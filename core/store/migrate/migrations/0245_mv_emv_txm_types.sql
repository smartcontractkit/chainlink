-- +goose Up
ALTER TYPE eth_txes_state RENAME TO txes_state;
ALTER TYPE txes_state SET SCHEMA "evm";

ALTER TYPE eth_tx_attempts_state RENAME TO tx_attempts_state;
ALTER TYPE tx_attempts_state SET SCHEMA "evm";
-- +goose Down
ALTER TYPE evm.txes_state SET SCHEMA "public";
ALTER TYPE txes_state RENAME TO eth_txes_state;

ALTER TYPE evm.tx_attempts_state SET SCHEMA "public";
ALTER TYPE tx_attempts_state RENAME TO eth_tx_attempts_state;

