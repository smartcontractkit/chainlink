-- +goose Up
CREATE SCHEMA IF NOT EXISTS aptos;

CREATE TABLE IF NOT EXISTS aptos.events (
    id BIGSERIAL PRIMARY KEY,
    event_account_address TEXT NOT NULL,
    event_handle TEXT NOT NULL,
    event_field_name TEXT NOT NULL,
    event_offset BIGINT NOT NULL,
    tx_version BIGINT NOT NULL,
    block_height TEXT NOT NULL,
    block_hash BYTEA NOT NULL,
    block_timestamp BIGINT NOT NULL,
    data JSONB NOT NULL,
    UNIQUE (event_account_address, event_handle, event_field_name, event_offset, tx_version)
);

CREATE INDEX IF NOT EXISTS idx_events_account_handle_offset
ON aptos.events(event_account_address, event_handle, event_field_name, tx_version, event_offset);

CREATE TABLE IF NOT EXISTS aptos.transmitter_sequence_nums (
    id BIGSERIAL PRIMARY KEY,
    transmitter_address TEXT NOT NULL,
    sequence_number BIGINT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (transmitter_address)
);

-- +goose Down
DROP TABLE IF EXISTS aptos.transmitter_sequence_nums;
DROP INDEX IF EXISTS idx_events_account_handle_offset;
DROP TABLE IF EXISTS aptos.events;
DROP SCHEMA IF EXISTS aptos;
