-- +goose Up

ALTER TABLE evm.log_poller_filters
    ADD COLUMN topic2 BYTEA CHECK (octet_length(topic2) = 32),
    ADD COLUMN topic3 BYTEA CHECK (octet_length(topic3) = 32),
    ADD COLUMN topic4 BYTEA CHECK (octet_length(topic4) = 32),
    ADD COLUMN max_logs_kept NUMERIC(78,0),
    ADD COLUMN logs_per_block NUMERIC(78,0),
    DROP CONSTRAINT evm_log_poller_filters_name_evm_chain_id_address_event_key;
-- Ordinary UNIQUE CONSTRAINT can't work for topics because they can be NULL. Any row with any column being NULL automatically satisfies the unique constraint (NULL != NULL)
-- Using a hash of all the columns treats NULL's as the same as any other field. If we ever get to a point where we can require postgresql >= 15 then this can
-- be fixed by using UNIQUE CONSTRAINT NULLS NOT DISTINCT which treats NULL's as if they were ordinary values (NULL == NULL)
CREATE UNIQUE INDEX evm_log_poller_filters_name_chain_address_event_topics_key ON evm.log_poller_filters (hash_record_extended((name, evm_chain_id, address, event, topic2, topic3, topic4), 0));

-- +goose Down

DROP INDEX IF EXISTS evm_log_poller_filters_name_chain_address_event_topics_key;
ALTER TABLE evm.log_poller_filters
    ADD CONSTRAINT evm_log_poller_filters_name_evm_chain_id_address_event_key UNIQUE (name, evm_chain_id, address, event),
    DROP COLUMN topic2,
    DROP COLUMN topic3,
    DROP COLUMN topic4,
    DROP COLUMN max_logs_kept,
    DROP COLUMN logs_per_block;
