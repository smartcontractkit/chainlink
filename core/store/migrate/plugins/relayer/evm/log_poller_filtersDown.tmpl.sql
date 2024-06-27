INSERT INTO evm.log_poller_filters ("name",address,"event", evm_chain_id, created_at,retention,topic2,topic3,topic4,max_logs_kept,logs_per_block)
SELECT "name", address, "event", '{{ .ChainID }}', created_at, retention, topic2, topic3, topic4, max_logs_kept, logs_per_block FROM {{ .Schema }}.log_poller_filters;

DROP TABLE {{ .Schema }}.log_poller_filters;
DROP FUNCTION evm.f_log_poller_filter_hash_no_chain;
