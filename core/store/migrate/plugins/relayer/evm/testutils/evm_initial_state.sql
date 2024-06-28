-- This is a snapshot of the `evm` schema after core node migration version 0244.
-- It is used to test plugin template migrations.

CREATE SCHEMA IF NOT EXISTS evm;

-- TODO: remove this after fixed in core 
CREATE TYPE public.eth_txes_state AS ENUM (
    'unstarted',
    'in_progress',
    'fatal_error',
    'unconfirmed',
    'confirmed_missing_receipt',
    'confirmed'
);
CREATE TYPE public.eth_tx_attempts_state AS ENUM (
    'in_progress',
    'insufficient_eth',
    'broadcast'
);



-- evm.forwarders definition

-- Drop table

-- DROP TABLE evm.forwarders;

CREATE TABLE evm.forwarders (
	id bigserial NOT NULL,
	address bytea NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	evm_chain_id numeric(78) NOT NULL,
	CONSTRAINT chk_address_length CHECK ((octet_length(address) = 20)),
	CONSTRAINT evm_forwarders_address_key UNIQUE (address),
	CONSTRAINT evm_forwarders_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_forwarders_created_at ON evm.forwarders USING brin (created_at);
CREATE INDEX idx_forwarders_evm_address ON evm.forwarders USING btree (address);
CREATE INDEX idx_forwarders_evm_chain_id ON evm.forwarders USING btree (evm_chain_id);
CREATE INDEX idx_forwarders_updated_at ON evm.forwarders USING brin (updated_at);


-- evm.heads definition

-- Drop table

-- DROP TABLE evm.heads;

CREATE TABLE evm.heads (
	id bigserial NOT NULL,
	hash bytea NOT NULL,
	"number" int8 NOT NULL,
	parent_hash bytea NOT NULL,
	created_at timestamptz NOT NULL,
	"timestamp" timestamptz NOT NULL,
	l1_block_number int8 NULL,
	evm_chain_id numeric(78) NOT NULL,
	base_fee_per_gas numeric(78) NULL,
	CONSTRAINT chk_hash_size CHECK ((octet_length(hash) = 32)),
	CONSTRAINT chk_parent_hash_size CHECK ((octet_length(parent_hash) = 32)),
	CONSTRAINT heads_pkey1 PRIMARY KEY (id)
);
CREATE UNIQUE INDEX idx_heads_evm_chain_id_hash ON evm.heads USING btree (evm_chain_id, hash);
CREATE INDEX idx_heads_evm_chain_id_number ON evm.heads USING btree (evm_chain_id, number);


-- evm.key_states definition

-- Drop table

-- DROP TABLE evm.key_states;

CREATE TABLE evm.key_states (
	id serial4 NOT NULL,
	address bytea NOT NULL,
	disabled bool DEFAULT false NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	evm_chain_id numeric(78) NOT NULL,
	CONSTRAINT chk_address_length CHECK ((octet_length(address) = 20)),
	CONSTRAINT eth_key_states_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_evm_key_states_address ON evm.key_states USING btree (address);
CREATE UNIQUE INDEX idx_evm_key_states_evm_chain_id_address ON evm.key_states USING btree (evm_chain_id, address);


-- evm.log_poller_blocks definition

-- Drop table

-- DROP TABLE evm.log_poller_blocks;

CREATE TABLE evm.log_poller_blocks (
	evm_chain_id numeric(78) NOT NULL,
	block_hash bytea NOT NULL,
	block_number int8 NOT NULL,
	created_at timestamptz NOT NULL,
	block_timestamp timestamptz NOT NULL,
	finalized_block_number int8 DEFAULT 0 NOT NULL,
	CONSTRAINT block_hash_uniq UNIQUE (evm_chain_id, block_hash),
	CONSTRAINT log_poller_blocks_block_number_check CHECK ((block_number > 0)),
	CONSTRAINT log_poller_blocks_finalized_block_number_check CHECK ((finalized_block_number >= 0)),
	CONSTRAINT log_poller_blocks_pkey PRIMARY KEY (block_number, evm_chain_id)
);
CREATE INDEX idx_evm_log_poller_blocks_order_by_block ON evm.log_poller_blocks USING btree (evm_chain_id, block_number DESC);


-- evm.log_poller_filters definition

-- Drop table

-- DROP TABLE evm.log_poller_filters;

-- DROP FUNCTION evm.f_log_poller_filter_hash(text, numeric, bytea, bytea, bytea, bytea, bytea);
CREATE OR REPLACE FUNCTION evm.f_log_poller_filter_hash(name text, evm_chain_id numeric, address bytea, event bytea, topic2 bytea, topic3 bytea, topic4 bytea)
 RETURNS bigint
 LANGUAGE sql
 IMMUTABLE PARALLEL SAFE COST 25
AS $function$SELECT hashtextextended(textin(record_out(($1,$2,$3,$4,$5,$6,$7))), 0)$function$
;

CREATE TABLE evm.log_poller_filters (
	id bigserial NOT NULL,
	"name" text NOT NULL,
	address bytea NOT NULL,
	"event" bytea NOT NULL,
	evm_chain_id numeric(78) NULL,
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
CREATE UNIQUE INDEX log_poller_filters_hash_key ON evm.log_poller_filters USING btree (evm.f_log_poller_filter_hash(name, evm_chain_id, address, event, topic2, topic3, topic4));





-- evm.logs definition

-- Drop table

-- DROP TABLE evm.logs;

CREATE TABLE evm.logs (
	evm_chain_id numeric(78) NOT NULL,
	log_index int8 NOT NULL,
	block_hash bytea NOT NULL,
	block_number int8 NOT NULL,
	address bytea NOT NULL,
	event_sig bytea NOT NULL,
	topics _bytea NOT NULL,
	tx_hash bytea NOT NULL,
	"data" bytea NOT NULL,
	created_at timestamptz NOT NULL,
	block_timestamp timestamptz NOT NULL,
	CONSTRAINT logs_block_number_check CHECK ((block_number > 0)),
	CONSTRAINT logs_pkey PRIMARY KEY (block_hash, log_index, evm_chain_id)
);
CREATE INDEX evm_logs_by_timestamp ON evm.logs USING btree (evm_chain_id, address, event_sig, block_timestamp, block_number);
CREATE INDEX evm_logs_idx ON evm.logs USING btree (evm_chain_id, block_number, address, event_sig);
CREATE INDEX evm_logs_idx_data_word_five ON evm.logs USING btree (address, event_sig, evm_chain_id, "substring"(data, 129, 32));
CREATE INDEX evm_logs_idx_data_word_four ON evm.logs USING btree (SUBSTRING(data FROM 97 FOR 32));
CREATE INDEX evm_logs_idx_data_word_one ON evm.logs USING btree (SUBSTRING(data FROM 1 FOR 32));
CREATE INDEX evm_logs_idx_data_word_three ON evm.logs USING btree (SUBSTRING(data FROM 65 FOR 32));
CREATE INDEX evm_logs_idx_data_word_two ON evm.logs USING btree (SUBSTRING(data FROM 33 FOR 32));
CREATE INDEX evm_logs_idx_topic_four ON evm.logs USING btree ((topics[4]));
CREATE INDEX evm_logs_idx_topic_three ON evm.logs USING btree ((topics[3]));
CREATE INDEX evm_logs_idx_topic_two ON evm.logs USING btree ((topics[2]));
CREATE INDEX evm_logs_idx_tx_hash ON evm.logs USING btree (tx_hash);
CREATE INDEX idx_evm_logs_ordered_by_block_and_created_at ON evm.logs USING btree (evm_chain_id, address, event_sig, block_number, created_at);


-- evm.txes definition

-- Drop table

-- DROP TABLE evm.txes;

CREATE TABLE evm.txes (
	id bigserial NOT NULL,
	nonce int8 NULL,
	from_address bytea NOT NULL,
	to_address bytea NOT NULL,
	encoded_payload bytea NOT NULL,
	value numeric(78) NOT NULL,
	gas_limit int8 NOT NULL,
	error text NULL,
	broadcast_at timestamptz NULL,
	created_at timestamptz NOT NULL,
	state public."eth_txes_state" DEFAULT 'unstarted'::eth_txes_state NOT NULL,
	meta jsonb NULL,
	subject uuid NULL,
	pipeline_task_run_id uuid NULL,
	min_confirmations int4 NULL,
	evm_chain_id numeric(78) NOT NULL,
	transmit_checker jsonb NULL,
	initial_broadcast_at timestamptz NULL,
	idempotency_key varchar(2000) NULL,
	signal_callback bool DEFAULT false NULL,
	callback_completed bool DEFAULT false NULL,
	CONSTRAINT chk_broadcast_at_is_sane CHECK ((broadcast_at > '2018-12-31 17:00:00-07'::timestamp with time zone)),
	CONSTRAINT chk_error_cannot_be_empty CHECK (((error IS NULL) OR (length(error) > 0))),
	CONSTRAINT chk_eth_txes_fsm CHECK ((((state = 'unstarted'::eth_txes_state) AND (nonce IS NULL) AND (error IS NULL) AND (broadcast_at IS NULL) AND (initial_broadcast_at IS NULL)) OR ((state = 'in_progress'::eth_txes_state) AND (nonce IS NOT NULL) AND (error IS NULL) AND (broadcast_at IS NULL) AND (initial_broadcast_at IS NULL)) OR ((state = 'fatal_error'::eth_txes_state) AND (error IS NOT NULL)) OR ((state = 'unconfirmed'::eth_txes_state) AND (nonce IS NOT NULL) AND (error IS NULL) AND (broadcast_at IS NOT NULL) AND (initial_broadcast_at IS NOT NULL)) OR ((state = 'confirmed'::eth_txes_state) AND (nonce IS NOT NULL) AND (error IS NULL) AND (broadcast_at IS NOT NULL) AND (initial_broadcast_at IS NOT NULL)) OR ((state = 'confirmed_missing_receipt'::eth_txes_state) AND (nonce IS NOT NULL) AND (error IS NULL) AND (broadcast_at IS NOT NULL) AND (initial_broadcast_at IS NOT NULL)))) NOT VALID,
	CONSTRAINT chk_from_address_length CHECK ((octet_length(from_address) = 20)),
	CONSTRAINT chk_to_address_length CHECK ((octet_length(to_address) = 20)),
	CONSTRAINT eth_txes_idempotency_key_key UNIQUE (idempotency_key),
	CONSTRAINT eth_txes_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_eth_txes_broadcast_at ON evm.txes USING brin (broadcast_at);
CREATE INDEX idx_eth_txes_created_at ON evm.txes USING brin (created_at);
CREATE INDEX idx_eth_txes_from_address ON evm.txes USING btree (from_address);
CREATE INDEX idx_eth_txes_initial_broadcast_at ON evm.txes USING brin (initial_broadcast_at);
CREATE INDEX idx_eth_txes_min_unconfirmed_nonce_for_key_evm_chain_id ON evm.txes USING btree (evm_chain_id, from_address, nonce) WHERE (state = 'unconfirmed'::eth_txes_state);
CREATE UNIQUE INDEX idx_eth_txes_nonce_from_address_per_evm_chain_id ON evm.txes USING btree (evm_chain_id, from_address, nonce);
CREATE UNIQUE INDEX idx_eth_txes_pipeline_run_task_id ON evm.txes USING btree (pipeline_task_run_id) WHERE (pipeline_task_run_id IS NOT NULL);
CREATE INDEX idx_eth_txes_state_from_address_evm_chain_id ON evm.txes USING btree (evm_chain_id, from_address, state) WHERE (state <> 'confirmed'::eth_txes_state);
CREATE INDEX idx_eth_txes_unstarted_subject_id_evm_chain_id ON evm.txes USING btree (evm_chain_id, subject, id) WHERE ((subject IS NOT NULL) AND (state = 'unstarted'::eth_txes_state));
CREATE UNIQUE INDEX idx_only_one_in_progress_tx_per_account_id_per_evm_chain_id ON evm.txes USING btree (evm_chain_id, from_address) WHERE (state = 'in_progress'::eth_txes_state);


-- evm.upkeep_states definition

-- Drop table

-- DROP TABLE evm.upkeep_states;

CREATE TABLE evm.upkeep_states (
	id serial4 NOT NULL,
	work_id text NOT NULL,
	evm_chain_id numeric(20) NOT NULL,
	upkeep_id numeric(78) NOT NULL,
	completion_state int2 NOT NULL,
	ineligibility_reason int2 NOT NULL,
	block_number int8 NOT NULL,
	inserted_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	CONSTRAINT evm_upkeep_states_pkey PRIMARY KEY (id),
	CONSTRAINT work_id_len_chk CHECK (((length(work_id) > 0) AND (length(work_id) < 255)))
);
CREATE INDEX idx_evm_upkeep_state_added_at_chain_id ON evm.upkeep_states USING btree (evm_chain_id, inserted_at);
CREATE UNIQUE INDEX idx_evm_upkeep_state_chainid_workid ON evm.upkeep_states USING btree (evm_chain_id, work_id);


-- evm.tx_attempts definition

-- Drop table

-- DROP TABLE evm.tx_attempts;

CREATE TABLE evm.tx_attempts (
	id bigserial NOT NULL,
	eth_tx_id int8 NOT NULL,
	gas_price numeric(78) NULL,
	signed_raw_tx bytea NOT NULL,
	hash bytea NOT NULL,
	broadcast_before_block_num int8 NULL,
	state public."eth_tx_attempts_state" NOT NULL,
	created_at timestamptz NOT NULL,
	chain_specific_gas_limit int8 NOT NULL,
	tx_type int2 DEFAULT 0 NOT NULL,
	gas_tip_cap numeric(78) NULL,
	gas_fee_cap numeric(78) NULL,
	is_purge_attempt bool DEFAULT false NOT NULL,
	CONSTRAINT chk_cannot_broadcast_before_block_zero CHECK (((broadcast_before_block_num IS NULL) OR (broadcast_before_block_num > 0))),
	CONSTRAINT chk_chain_specific_gas_limit_not_zero CHECK ((chain_specific_gas_limit > 0)),
	CONSTRAINT chk_eth_tx_attempts_fsm CHECK ((((state = ANY (ARRAY['in_progress'::eth_tx_attempts_state, 'insufficient_eth'::eth_tx_attempts_state])) AND (broadcast_before_block_num IS NULL)) OR (state = 'broadcast'::eth_tx_attempts_state))),
	CONSTRAINT chk_hash_length CHECK ((octet_length(hash) = 32)),
	CONSTRAINT chk_legacy_or_dynamic CHECK ((((tx_type = 0) AND (gas_price IS NOT NULL) AND (gas_tip_cap IS NULL) AND (gas_fee_cap IS NULL)) OR ((tx_type = 2) AND (gas_price IS NULL) AND (gas_tip_cap IS NOT NULL) AND (gas_fee_cap IS NOT NULL)))),
	CONSTRAINT chk_sanity_fee_cap_tip_cap CHECK (((gas_tip_cap IS NULL) OR (gas_fee_cap IS NULL) OR (gas_tip_cap <= gas_fee_cap))),
	CONSTRAINT chk_signed_raw_tx_present CHECK ((octet_length(signed_raw_tx) > 0)),
	CONSTRAINT chk_tx_type_is_byte CHECK (((tx_type >= 0) AND (tx_type <= 255))),
	CONSTRAINT eth_tx_attempts_pkey PRIMARY KEY (id),
	CONSTRAINT eth_tx_attempts_eth_tx_id_fkey FOREIGN KEY (eth_tx_id) REFERENCES evm.txes(id) ON DELETE CASCADE
);
CREATE INDEX idx_eth_tx_attempts_broadcast_before_block_num ON evm.tx_attempts USING btree (broadcast_before_block_num);
CREATE INDEX idx_eth_tx_attempts_created_at ON evm.tx_attempts USING brin (created_at);
CREATE UNIQUE INDEX idx_eth_tx_attempts_hash ON evm.tx_attempts USING btree (hash);
CREATE INDEX idx_eth_tx_attempts_unbroadcast ON evm.tx_attempts USING btree (state) WHERE (state <> 'broadcast'::eth_tx_attempts_state);
CREATE UNIQUE INDEX idx_eth_tx_attempts_unique_gas_prices ON evm.tx_attempts USING btree (eth_tx_id, gas_price);
CREATE UNIQUE INDEX idx_only_one_unbroadcast_attempt_per_eth_tx ON evm.tx_attempts USING btree (eth_tx_id) WHERE (state <> 'broadcast'::eth_tx_attempts_state);


-- evm.receipts definition

-- Drop table

-- DROP TABLE evm.receipts;

CREATE TABLE evm.receipts (
	id bigserial NOT NULL,
	tx_hash bytea NOT NULL,
	block_hash bytea NOT NULL,
	block_number int8 NOT NULL,
	transaction_index int8 NOT NULL,
	receipt jsonb NOT NULL,
	created_at timestamptz NOT NULL,
	CONSTRAINT chk_hash_length CHECK (((octet_length(tx_hash) = 32) AND (octet_length(block_hash) = 32))),
	CONSTRAINT eth_receipts_pkey PRIMARY KEY (id),
	CONSTRAINT eth_receipts_tx_hash_fkey FOREIGN KEY (tx_hash) REFERENCES evm.tx_attempts(hash) ON DELETE CASCADE
);
CREATE INDEX idx_eth_receipts_block_number ON evm.receipts USING btree (block_number);
CREATE INDEX idx_eth_receipts_created_at ON evm.receipts USING brin (created_at);
CREATE UNIQUE INDEX idx_eth_receipts_unique ON evm.receipts USING btree (tx_hash, block_hash);