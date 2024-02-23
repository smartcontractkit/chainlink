-- +goose Up

-- This generates a unique BIGINT for the log_poller_filters table from hashing (name, evm_chain_id, address, event, topic2, topic3, topic4).
-- An ordinary UNIQUE CONSTRAINT can't work for this because the topics can be NULL. Any row with any column being NULL automatically satisfies the unique constraint (NULL != NULL)
-- There are simpler ways of doing this in postgres 12, 13, 14, and especially 15. But for now, we still officially support postgresql 11 and this should be just as efficient.
CREATE OR REPLACE FUNCTION evm.f_log_filter_row_id(name TEXT, evm_chain_id NUMERIC, address BYTEA, event BYTEA, topic2 BYTEA, topic3 BYTEA, topic4 BYTEA)
   RETURNS BIGINT
   LANGUAGE SQL IMMUTABLE COST 25 PARALLEL SAFE AS 'SELECT 2^32 * hashtext(textin(record_out(($1,$3,$5,$7)))) + hashtext(textin(record_out(($2, $4, $6))))';

ALTER TABLE evm.log_poller_filters
    ADD COLUMN topic2 BYTEA CHECK (octet_length(topic2) = 32),
    ADD COLUMN topic3 BYTEA CHECK (octet_length(topic3) = 32),
    ADD COLUMN topic4 BYTEA CHECK (octet_length(topic4) = 32),
    ADD COLUMN max_logs_kept BIGINT,
    ADD COLUMN logs_per_block BIGINT;

CREATE UNIQUE INDEX evm_log_poller_filters_name_chain_address_event_topics_key ON evm.log_poller_filters (evm.f_log_filter_row_id(name, evm_chain_id, address, event, topic2, topic3, topic4));

ALTER TABLE evm.log_poller_filters
    DROP CONSTRAINT evm_log_poller_filters_name_evm_chain_id_address_event_key;

-- +goose Down

ALTER TABLE evm.log_poller_filters
    ADD CONSTRAINT evm_log_poller_filters_name_evm_chain_id_address_event_key UNIQUE (name, evm_chain_id, address, event);
DROP INDEX IF EXISTS evm_log_poller_filters_name_chain_address_event_topics_key;

ALTER TABLE evm.log_poller_filters
    DROP COLUMN topic2,
    DROP COLUMN topic3,
    DROP COLUMN topic4,
    DROP COLUMN max_logs_kept,
    DROP COLUMN logs_per_block;

DROP FUNCTION IF EXISTS evm.f_log_filter_row_id(TEXT, NUMERIC, BYTEA, BYTEA, BYTEA, BYTEA, BYTEA);
