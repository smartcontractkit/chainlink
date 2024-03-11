-- +goose Up

-- This generates a unique BIGINT for the log_poller_filters table from hashing (name, evm_chain_id, address, event, topic2, topic3, topic4).
-- Using an ordinary multi-column index on 7 columns would require a lot of extra storage space, and there are additional complications due to the topics being allowed to be NULL.
-- Note for updating this if and when we drop support for postgresql 12 & 13: hashrecordextended() can be used directly in postgresql 14, avoiding the need for a helper function.
-- The helper function is necessary only for the IMMUTABLE keyword
CREATE OR REPLACE FUNCTION evm.f_log_poller_filter_hash(name TEXT, evm_chain_id NUMERIC, address BYTEA, event BYTEA, topic2 BYTEA, topic3 BYTEA, topic4 BYTEA)
   RETURNS BIGINT
   LANGUAGE SQL IMMUTABLE COST 25 PARALLEL SAFE AS 'SELECT hashtextextended(textin(record_out(($1,$2,$3,$4,$5,$6,$7))), 0)';

ALTER TABLE evm.log_poller_filters
    ADD COLUMN topic2 BYTEA CHECK (octet_length(topic2) = 32),
    ADD COLUMN topic3 BYTEA CHECK (octet_length(topic3) = 32),
    ADD COLUMN topic4 BYTEA CHECK (octet_length(topic4) = 32),
    ADD COLUMN max_logs_kept BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN logs_per_block BIGINT NOT NULL DEFAULT 0;

CREATE UNIQUE INDEX log_poller_filters_hash_key ON evm.log_poller_filters (evm.f_log_poller_filter_hash(name, evm_chain_id, address, event, topic2, topic3, topic4));

ALTER TABLE evm.log_poller_filters
    DROP CONSTRAINT evm_log_poller_filters_name_evm_chain_id_address_event_key;

-- +goose Down

ALTER TABLE evm.log_poller_filters
    ADD CONSTRAINT evm_log_poller_filters_name_evm_chain_id_address_event_key UNIQUE (name, evm_chain_id, address, event);
DROP INDEX IF EXISTS log_poller_filters_hash_key;

ALTER TABLE evm.log_poller_filters
    DROP COLUMN topic2,
    DROP COLUMN topic3,
    DROP COLUMN topic4,
    DROP COLUMN max_logs_kept,
    DROP COLUMN logs_per_block;

DROP FUNCTION IF EXISTS evm.f_log_poller_filter_hash(TEXT, NUMERIC, BYTEA, BYTEA, BYTEA, BYTEA, BYTEA);
