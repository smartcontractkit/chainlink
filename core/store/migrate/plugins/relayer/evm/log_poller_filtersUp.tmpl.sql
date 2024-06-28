CREATE TABLE {{ .Schema }}.log_poller_filters (
	id bigserial NOT NULL,
	"name" text NOT NULL,
	address bytea NOT NULL,
	"event" bytea NOT NULL,
--	evm_chain_id numeric(78) NULL,
	created_at timestamptz NOT NULL,
	retention int8 DEFAULT 0 NULL,
	topic2 bytea NULL,
	topic3 bytea NULL,
	topic4 bytea NULL,
	max_logs_kept int8 DEFAULT 0 NOT NULL,
	logs_per_block int8 DEFAULT 0 NOT NULL,
	CONSTRAINT evm_log_poller_filters_address_check CHECK ((octet_length(address) = 20)),
	CONSTRAINT evm_log_poller_filters_event_check CHECK ((octet_length(event) = 32)),
	CONSTRAINT evm_log_poller_filters_name_check CHECK ((length(name) > 0)),
	CONSTRAINT evm_log_poller_filters_pkey PRIMARY KEY (id),
	CONSTRAINT log_poller_filters_topic2_check CHECK ((octet_length(topic2) = 32)),
	CONSTRAINT log_poller_filters_topic3_check CHECK ((octet_length(topic3) = 32)),
	CONSTRAINT log_poller_filters_topic4_check CHECK ((octet_length(topic4) = 32))
);
-- what to do with the log poller function. doesn't make sense to copy it over to the new schema

CREATE OR REPLACE FUNCTION evm.f_log_poller_filter_hash_no_chain(name text, address bytea, event bytea, topic2 bytea, topic3 bytea, topic4 bytea)
 RETURNS bigint
 LANGUAGE sql
 IMMUTABLE PARALLEL SAFE COST 25
AS $function$SELECT hashtextextended(textin(record_out(($1,$2,$3,$4,$5,$6))), 0)$function$
;

CREATE UNIQUE INDEX log_poller_filters_hash_key ON {{ .Schema }}.log_poller_filters USING btree (evm.f_log_poller_filter_hash_no_chain(name,  address, event, topic2, topic3, topic4));

INSERT INTO {{ .Schema }}.log_poller_filters ("name",address,"event",created_at,retention,topic2,topic3,topic4,max_logs_kept,logs_per_block)
SELECT "name", address, "event", created_at, retention, topic2, topic3, topic4, max_logs_kept, logs_per_block FROM evm.log_poller_filters WHERE evm_chain_id = '{{ .ChainID}}';

DELETE FROM evm.log_poller_filters WHERE evm_chain_id = '{{ .ChainID}}';
