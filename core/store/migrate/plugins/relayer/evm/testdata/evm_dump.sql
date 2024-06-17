--
-- PostgreSQL database dump
--

-- Dumped from database version 15.7
-- Dumped by pg_dump version 15.7

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: evm; Type: SCHEMA; Schema: -; Owner: chainlink_dev
--

CREATE SCHEMA evm;


ALTER SCHEMA evm OWNER TO chainlink_dev;

--
-- Name: f_log_poller_filter_hash(text, numeric, bytea, bytea, bytea, bytea, bytea); Type: FUNCTION; Schema: evm; Owner: chainlink_dev
--

CREATE FUNCTION evm.f_log_poller_filter_hash(name text, evm_chain_id numeric, address bytea, event bytea, topic2 bytea, topic3 bytea, topic4 bytea) RETURNS bigint
    LANGUAGE sql IMMUTABLE COST 25 PARALLEL SAFE
    AS $_$SELECT hashtextextended(textin(record_out(($1,$2,$3,$4,$5,$6,$7))), 0)$_$;


ALTER FUNCTION evm.f_log_poller_filter_hash(name text, evm_chain_id numeric, address bytea, event bytea, topic2 bytea, topic3 bytea, topic4 bytea) OWNER TO chainlink_dev;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: key_states; Type: TABLE; Schema: evm; Owner: chainlink_dev
--

CREATE TABLE evm.key_states (
    id integer NOT NULL,
    address bytea NOT NULL,
    disabled boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    evm_chain_id numeric(78,0) NOT NULL,
    CONSTRAINT chk_address_length CHECK ((octet_length(address) = 20))
);


ALTER TABLE evm.key_states OWNER TO chainlink_dev;

--
-- Name: eth_key_states_id_seq; Type: SEQUENCE; Schema: evm; Owner: chainlink_dev
--

CREATE SEQUENCE evm.eth_key_states_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE evm.eth_key_states_id_seq OWNER TO chainlink_dev;

--
-- Name: eth_key_states_id_seq; Type: SEQUENCE OWNED BY; Schema: evm; Owner: chainlink_dev
--

ALTER SEQUENCE evm.eth_key_states_id_seq OWNED BY evm.key_states.id;


--
-- Name: receipts; Type: TABLE; Schema: evm; Owner: chainlink_dev
--

CREATE TABLE evm.receipts (
    id bigint NOT NULL,
    tx_hash bytea NOT NULL,
    block_hash bytea NOT NULL,
    block_number bigint NOT NULL,
    transaction_index bigint NOT NULL,
    receipt jsonb NOT NULL,
    created_at timestamp with time zone NOT NULL,
    CONSTRAINT chk_hash_length CHECK (((octet_length(tx_hash) = 32) AND (octet_length(block_hash) = 32)))
);


ALTER TABLE evm.receipts OWNER TO chainlink_dev;

--
-- Name: eth_receipts_id_seq; Type: SEQUENCE; Schema: evm; Owner: chainlink_dev
--

CREATE SEQUENCE evm.eth_receipts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE evm.eth_receipts_id_seq OWNER TO chainlink_dev;

--
-- Name: eth_receipts_id_seq; Type: SEQUENCE OWNED BY; Schema: evm; Owner: chainlink_dev
--

ALTER SEQUENCE evm.eth_receipts_id_seq OWNED BY evm.receipts.id;


--
-- Name: tx_attempts; Type: TABLE; Schema: evm; Owner: chainlink_dev
--

CREATE TABLE evm.tx_attempts (
    id bigint NOT NULL,
    eth_tx_id bigint NOT NULL,
    gas_price numeric(78,0),
    signed_raw_tx bytea NOT NULL,
    hash bytea NOT NULL,
    broadcast_before_block_num bigint,
    state public.eth_tx_attempts_state NOT NULL,
    created_at timestamp with time zone NOT NULL,
    chain_specific_gas_limit bigint NOT NULL,
    tx_type smallint DEFAULT 0 NOT NULL,
    gas_tip_cap numeric(78,0),
    gas_fee_cap numeric(78,0),
    is_purge_attempt boolean DEFAULT false NOT NULL,
    CONSTRAINT chk_cannot_broadcast_before_block_zero CHECK (((broadcast_before_block_num IS NULL) OR (broadcast_before_block_num > 0))),
    CONSTRAINT chk_chain_specific_gas_limit_not_zero CHECK ((chain_specific_gas_limit > 0)),
    CONSTRAINT chk_eth_tx_attempts_fsm CHECK ((((state = ANY (ARRAY['in_progress'::public.eth_tx_attempts_state, 'insufficient_eth'::public.eth_tx_attempts_state])) AND (broadcast_before_block_num IS NULL)) OR (state = 'broadcast'::public.eth_tx_attempts_state))),
    CONSTRAINT chk_hash_length CHECK ((octet_length(hash) = 32)),
    CONSTRAINT chk_legacy_or_dynamic CHECK ((((tx_type = 0) AND (gas_price IS NOT NULL) AND (gas_tip_cap IS NULL) AND (gas_fee_cap IS NULL)) OR ((tx_type = 2) AND (gas_price IS NULL) AND (gas_tip_cap IS NOT NULL) AND (gas_fee_cap IS NOT NULL)))),
    CONSTRAINT chk_sanity_fee_cap_tip_cap CHECK (((gas_tip_cap IS NULL) OR (gas_fee_cap IS NULL) OR (gas_tip_cap <= gas_fee_cap))),
    CONSTRAINT chk_signed_raw_tx_present CHECK ((octet_length(signed_raw_tx) > 0)),
    CONSTRAINT chk_tx_type_is_byte CHECK (((tx_type >= 0) AND (tx_type <= 255)))
);


ALTER TABLE evm.tx_attempts OWNER TO chainlink_dev;

--
-- Name: eth_tx_attempts_id_seq; Type: SEQUENCE; Schema: evm; Owner: chainlink_dev
--

CREATE SEQUENCE evm.eth_tx_attempts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE evm.eth_tx_attempts_id_seq OWNER TO chainlink_dev;

--
-- Name: eth_tx_attempts_id_seq; Type: SEQUENCE OWNED BY; Schema: evm; Owner: chainlink_dev
--

ALTER SEQUENCE evm.eth_tx_attempts_id_seq OWNED BY evm.tx_attempts.id;


--
-- Name: txes; Type: TABLE; Schema: evm; Owner: chainlink_dev
--

CREATE TABLE evm.txes (
    id bigint NOT NULL,
    nonce bigint,
    from_address bytea NOT NULL,
    to_address bytea NOT NULL,
    encoded_payload bytea NOT NULL,
    value numeric(78,0) NOT NULL,
    gas_limit bigint NOT NULL,
    error text,
    broadcast_at timestamp with time zone,
    created_at timestamp with time zone NOT NULL,
    state public.eth_txes_state DEFAULT 'unstarted'::public.eth_txes_state NOT NULL,
    meta jsonb,
    subject uuid,
    pipeline_task_run_id uuid,
    min_confirmations integer,
    evm_chain_id numeric(78,0) NOT NULL,
    transmit_checker jsonb,
    initial_broadcast_at timestamp with time zone,
    idempotency_key character varying(2000),
    signal_callback boolean DEFAULT false,
    callback_completed boolean DEFAULT false,
    CONSTRAINT chk_broadcast_at_is_sane CHECK ((broadcast_at > '2019-01-01 00:00:00+00'::timestamp with time zone)),
    CONSTRAINT chk_error_cannot_be_empty CHECK (((error IS NULL) OR (length(error) > 0))),
    CONSTRAINT chk_from_address_length CHECK ((octet_length(from_address) = 20)),
    CONSTRAINT chk_to_address_length CHECK ((octet_length(to_address) = 20))
);


ALTER TABLE evm.txes OWNER TO chainlink_dev;

--
-- Name: eth_txes_id_seq; Type: SEQUENCE; Schema: evm; Owner: chainlink_dev
--

CREATE SEQUENCE evm.eth_txes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE evm.eth_txes_id_seq OWNER TO chainlink_dev;

--
-- Name: eth_txes_id_seq; Type: SEQUENCE OWNED BY; Schema: evm; Owner: chainlink_dev
--

ALTER SEQUENCE evm.eth_txes_id_seq OWNED BY evm.txes.id;


--
-- Name: forwarders; Type: TABLE; Schema: evm; Owner: chainlink_dev
--

CREATE TABLE evm.forwarders (
    id bigint NOT NULL,
    address bytea NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    evm_chain_id numeric(78,0) NOT NULL,
    CONSTRAINT chk_address_length CHECK ((octet_length(address) = 20))
);


ALTER TABLE evm.forwarders OWNER TO chainlink_dev;

--
-- Name: evm_forwarders_id_seq; Type: SEQUENCE; Schema: evm; Owner: chainlink_dev
--

CREATE SEQUENCE evm.evm_forwarders_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE evm.evm_forwarders_id_seq OWNER TO chainlink_dev;

--
-- Name: evm_forwarders_id_seq; Type: SEQUENCE OWNED BY; Schema: evm; Owner: chainlink_dev
--

ALTER SEQUENCE evm.evm_forwarders_id_seq OWNED BY evm.forwarders.id;


--
-- Name: log_poller_filters; Type: TABLE; Schema: evm; Owner: chainlink_dev
--

CREATE TABLE evm.log_poller_filters (
    id bigint NOT NULL,
    name text NOT NULL,
    address bytea NOT NULL,
    event bytea NOT NULL,
    evm_chain_id numeric(78,0),
    created_at timestamp with time zone NOT NULL,
    retention bigint DEFAULT 0,
    topic2 bytea,
    topic3 bytea,
    topic4 bytea,
    max_logs_kept bigint DEFAULT 0 NOT NULL,
    logs_per_block bigint DEFAULT 0 NOT NULL,
    CONSTRAINT evm_log_poller_filters_address_check CHECK ((octet_length(address) = 20)),
    CONSTRAINT evm_log_poller_filters_event_check CHECK ((octet_length(event) = 32)),
    CONSTRAINT evm_log_poller_filters_name_check CHECK ((length(name) > 0)),
    CONSTRAINT log_poller_filters_topic2_check CHECK ((octet_length(topic2) = 32)),
    CONSTRAINT log_poller_filters_topic3_check CHECK ((octet_length(topic3) = 32)),
    CONSTRAINT log_poller_filters_topic4_check CHECK ((octet_length(topic4) = 32))
);


ALTER TABLE evm.log_poller_filters OWNER TO chainlink_dev;

--
-- Name: evm_log_poller_filters_id_seq; Type: SEQUENCE; Schema: evm; Owner: chainlink_dev
--

CREATE SEQUENCE evm.evm_log_poller_filters_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE evm.evm_log_poller_filters_id_seq OWNER TO chainlink_dev;

--
-- Name: evm_log_poller_filters_id_seq; Type: SEQUENCE OWNED BY; Schema: evm; Owner: chainlink_dev
--

ALTER SEQUENCE evm.evm_log_poller_filters_id_seq OWNED BY evm.log_poller_filters.id;


--
-- Name: upkeep_states; Type: TABLE; Schema: evm; Owner: chainlink_dev
--

CREATE TABLE evm.upkeep_states (
    id integer NOT NULL,
    work_id text NOT NULL,
    evm_chain_id numeric(20,0) NOT NULL,
    upkeep_id numeric(78,0) NOT NULL,
    completion_state smallint NOT NULL,
    ineligibility_reason smallint NOT NULL,
    block_number bigint NOT NULL,
    inserted_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT work_id_len_chk CHECK (((length(work_id) > 0) AND (length(work_id) < 255)))
);


ALTER TABLE evm.upkeep_states OWNER TO chainlink_dev;

--
-- Name: evm_upkeep_states_id_seq; Type: SEQUENCE; Schema: evm; Owner: chainlink_dev
--

CREATE SEQUENCE evm.evm_upkeep_states_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE evm.evm_upkeep_states_id_seq OWNER TO chainlink_dev;

--
-- Name: evm_upkeep_states_id_seq; Type: SEQUENCE OWNED BY; Schema: evm; Owner: chainlink_dev
--

ALTER SEQUENCE evm.evm_upkeep_states_id_seq OWNED BY evm.upkeep_states.id;


--
-- Name: heads; Type: TABLE; Schema: evm; Owner: chainlink_dev
--

CREATE TABLE evm.heads (
    id bigint NOT NULL,
    hash bytea NOT NULL,
    number bigint NOT NULL,
    parent_hash bytea NOT NULL,
    created_at timestamp with time zone NOT NULL,
    "timestamp" timestamp with time zone NOT NULL,
    l1_block_number bigint,
    evm_chain_id numeric(78,0) NOT NULL,
    base_fee_per_gas numeric(78,0),
    CONSTRAINT chk_hash_size CHECK ((octet_length(hash) = 32)),
    CONSTRAINT chk_parent_hash_size CHECK ((octet_length(parent_hash) = 32))
);


ALTER TABLE evm.heads OWNER TO chainlink_dev;

--
-- Name: heads_id_seq; Type: SEQUENCE; Schema: evm; Owner: chainlink_dev
--

CREATE SEQUENCE evm.heads_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE evm.heads_id_seq OWNER TO chainlink_dev;

--
-- Name: heads_id_seq; Type: SEQUENCE OWNED BY; Schema: evm; Owner: chainlink_dev
--

ALTER SEQUENCE evm.heads_id_seq OWNED BY evm.heads.id;


--
-- Name: log_poller_blocks; Type: TABLE; Schema: evm; Owner: chainlink_dev
--

CREATE TABLE evm.log_poller_blocks (
    evm_chain_id numeric(78,0) NOT NULL,
    block_hash bytea NOT NULL,
    block_number bigint NOT NULL,
    created_at timestamp with time zone NOT NULL,
    block_timestamp timestamp with time zone NOT NULL,
    finalized_block_number bigint DEFAULT 0 NOT NULL,
    CONSTRAINT log_poller_blocks_block_number_check CHECK ((block_number > 0)),
    CONSTRAINT log_poller_blocks_finalized_block_number_check CHECK ((finalized_block_number >= 0))
);


ALTER TABLE evm.log_poller_blocks OWNER TO chainlink_dev;

--
-- Name: logs; Type: TABLE; Schema: evm; Owner: chainlink_dev
--

CREATE TABLE evm.logs (
    evm_chain_id numeric(78,0) NOT NULL,
    log_index bigint NOT NULL,
    block_hash bytea NOT NULL,
    block_number bigint NOT NULL,
    address bytea NOT NULL,
    event_sig bytea NOT NULL,
    topics bytea[] NOT NULL,
    tx_hash bytea NOT NULL,
    data bytea NOT NULL,
    created_at timestamp with time zone NOT NULL,
    block_timestamp timestamp with time zone NOT NULL,
    CONSTRAINT logs_block_number_check CHECK ((block_number > 0))
);


ALTER TABLE evm.logs OWNER TO chainlink_dev;

--
-- Name: forwarders id; Type: DEFAULT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE ONLY evm.forwarders ALTER COLUMN id SET DEFAULT nextval('evm.evm_forwarders_id_seq'::regclass);


--
-- Name: heads id; Type: DEFAULT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE ONLY evm.heads ALTER COLUMN id SET DEFAULT nextval('evm.heads_id_seq'::regclass);


--
-- Name: key_states id; Type: DEFAULT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE ONLY evm.key_states ALTER COLUMN id SET DEFAULT nextval('evm.eth_key_states_id_seq'::regclass);


--
-- Name: log_poller_filters id; Type: DEFAULT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE ONLY evm.log_poller_filters ALTER COLUMN id SET DEFAULT nextval('evm.evm_log_poller_filters_id_seq'::regclass);


--
-- Name: receipts id; Type: DEFAULT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE ONLY evm.receipts ALTER COLUMN id SET DEFAULT nextval('evm.eth_receipts_id_seq'::regclass);


--
-- Name: tx_attempts id; Type: DEFAULT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE ONLY evm.tx_attempts ALTER COLUMN id SET DEFAULT nextval('evm.eth_tx_attempts_id_seq'::regclass);


--
-- Name: txes id; Type: DEFAULT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE ONLY evm.txes ALTER COLUMN id SET DEFAULT nextval('evm.eth_txes_id_seq'::regclass);


--
-- Name: upkeep_states id; Type: DEFAULT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE ONLY evm.upkeep_states ALTER COLUMN id SET DEFAULT nextval('evm.evm_upkeep_states_id_seq'::regclass);


--
-- Data for Name: forwarders; Type: TABLE DATA; Schema: evm; Owner: chainlink_dev
--

COPY evm.forwarders (id, address, created_at, updated_at, evm_chain_id) FROM stdin;
5       \\x3031323334353637383930313233343536373839     2024-06-17 16:56:35.302744+00   2024-06-17 16:56:35.302744+00   42
6       \\x3030303030303030303030303030303030303030     2024-06-17 16:57:30.232174+00   2024-06-17 16:57:30.232174+00   42
\.


--
-- Data for Name: heads; Type: TABLE DATA; Schema: evm; Owner: chainlink_dev
--

COPY evm.heads (id, hash, number, parent_hash, created_at, "timestamp", l1_block_number, evm_chain_id, base_fee_per_gas) FROM stdin;
\.


--
-- Data for Name: key_states; Type: TABLE DATA; Schema: evm; Owner: chainlink_dev
--

COPY evm.key_states (id, address, disabled, created_at, updated_at, evm_chain_id) FROM stdin;
\.


--
-- Data for Name: log_poller_blocks; Type: TABLE DATA; Schema: evm; Owner: chainlink_dev
--

COPY evm.log_poller_blocks (evm_chain_id, block_hash, block_number, created_at, block_timestamp, finalized_block_number) FROM stdin;
\.


--
-- Data for Name: log_poller_filters; Type: TABLE DATA; Schema: evm; Owner: chainlink_dev
--

COPY evm.log_poller_filters (id, name, address, event, evm_chain_id, created_at, retention, topic2, topic3, topic4, max_logs_kept, logs_per_block) FROM stdin;
\.


--
-- Data for Name: logs; Type: TABLE DATA; Schema: evm; Owner: chainlink_dev
--

COPY evm.logs (evm_chain_id, log_index, block_hash, block_number, address, event_sig, topics, tx_hash, data, created_at, block_timestamp) FROM stdin;
\.


--
-- Data for Name: receipts; Type: TABLE DATA; Schema: evm; Owner: chainlink_dev
--

COPY evm.receipts (id, tx_hash, block_hash, block_number, transaction_index, receipt, created_at) FROM stdin;
\.


--
-- Data for Name: tx_attempts; Type: TABLE DATA; Schema: evm; Owner: chainlink_dev
--

COPY evm.tx_attempts (id, eth_tx_id, gas_price, signed_raw_tx, hash, broadcast_before_block_num, state, created_at, chain_specific_gas_limit, tx_type, gas_tip_cap, gas_fee_cap, is_purge_attempt) FROM stdin;
\.


--
-- Data for Name: txes; Type: TABLE DATA; Schema: evm; Owner: chainlink_dev
--

COPY evm.txes (id, nonce, from_address, to_address, encoded_payload, value, gas_limit, error, broadcast_at, created_at, state, meta, subject, pipeline_task_run_id, min_confirmations, evm_chain_id, transmit_checker, initial_broadcast_at, idempotency_key, signal_callback, callback_completed) FROM stdin;
\.


--
-- Data for Name: upkeep_states; Type: TABLE DATA; Schema: evm; Owner: chainlink_dev
--

COPY evm.upkeep_states (id, work_id, evm_chain_id, upkeep_id, completion_state, ineligibility_reason, block_number, inserted_at) FROM stdin;
\.


--
-- Name: eth_key_states_id_seq; Type: SEQUENCE SET; Schema: evm; Owner: chainlink_dev
--

SELECT pg_catalog.setval('evm.eth_key_states_id_seq', 1, false);


--
-- Name: eth_receipts_id_seq; Type: SEQUENCE SET; Schema: evm; Owner: chainlink_dev
--

SELECT pg_catalog.setval('evm.eth_receipts_id_seq', 1, false);


--
-- Name: eth_tx_attempts_id_seq; Type: SEQUENCE SET; Schema: evm; Owner: chainlink_dev
--

SELECT pg_catalog.setval('evm.eth_tx_attempts_id_seq', 1, false);


--
-- Name: eth_txes_id_seq; Type: SEQUENCE SET; Schema: evm; Owner: chainlink_dev
--

SELECT pg_catalog.setval('evm.eth_txes_id_seq', 1, false);


--
-- Name: evm_forwarders_id_seq; Type: SEQUENCE SET; Schema: evm; Owner: chainlink_dev
--

SELECT pg_catalog.setval('evm.evm_forwarders_id_seq', 6, true);


--
-- Name: evm_log_poller_filters_id_seq; Type: SEQUENCE SET; Schema: evm; Owner: chainlink_dev
--

SELECT pg_catalog.setval('evm.evm_log_poller_filters_id_seq', 1, false);


--
-- Name: evm_upkeep_states_id_seq; Type: SEQUENCE SET; Schema: evm; Owner: chainlink_dev
--

SELECT pg_catalog.setval('evm.evm_upkeep_states_id_seq', 1, false);


--
-- Name: heads_id_seq; Type: SEQUENCE SET; Schema: evm; Owner: chainlink_dev
--

SELECT pg_catalog.setval('evm.heads_id_seq', 1, false);


--
-- Name: log_poller_blocks block_hash_uniq; Type: CONSTRAINT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE ONLY evm.log_poller_blocks
    ADD CONSTRAINT block_hash_uniq UNIQUE (evm_chain_id, block_hash);


--
-- Name: txes chk_eth_txes_fsm; Type: CHECK CONSTRAINT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE evm.txes
    ADD CONSTRAINT chk_eth_txes_fsm CHECK ((((state = 'unstarted'::public.eth_txes_state) AND (nonce IS NULL) AND (error IS NULL) AND (broadcast_at IS NULL) AND (initial_broadcast_at IS NULL)) OR ((state = 'in_progress'::public.eth_txes_state) AND (nonce IS NOT NULL) AND (error IS NULL) AND (broadcast_at IS NULL) AND (initial_broadcast_at IS NULL)) OR ((state = 'fatal_error'::public.eth_txes_state) AND (error IS NOT NULL)) OR ((state = 'unconfirmed'::public.eth_txes_state) AND (nonce IS NOT NULL) AND (error IS NULL) AND (broadcast_at IS NOT NULL) AND (initial_broadcast_at IS NOT NULL)) OR ((state = 'confirmed'::public.eth_txes_state) AND (nonce IS NOT NULL) AND (error IS NULL) AND (broadcast_at IS NOT NULL) AND (initial_broadcast_at IS NOT NULL)) OR ((state = 'confirmed_missing_receipt'::public.eth_txes_state) AND (nonce IS NOT NULL) AND (error IS NULL) AND (broadcast_at IS NOT NULL) AND (initial_broadcast_at IS NOT NULL)))) NOT VALID;


--
-- Name: key_states eth_key_states_pkey; Type: CONSTRAINT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE ONLY evm.key_states
    ADD CONSTRAINT eth_key_states_pkey PRIMARY KEY (id);


--
-- Name: receipts eth_receipts_pkey; Type: CONSTRAINT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE ONLY evm.receipts
    ADD CONSTRAINT eth_receipts_pkey PRIMARY KEY (id);


--
-- Name: tx_attempts eth_tx_attempts_pkey; Type: CONSTRAINT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE ONLY evm.tx_attempts
    ADD CONSTRAINT eth_tx_attempts_pkey PRIMARY KEY (id);


--
-- Name: txes eth_txes_idempotency_key_key; Type: CONSTRAINT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE ONLY evm.txes
    ADD CONSTRAINT eth_txes_idempotency_key_key UNIQUE (idempotency_key);


--
-- Name: txes eth_txes_pkey; Type: CONSTRAINT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE ONLY evm.txes
    ADD CONSTRAINT eth_txes_pkey PRIMARY KEY (id);


--
-- Name: forwarders evm_forwarders_address_key; Type: CONSTRAINT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE ONLY evm.forwarders
    ADD CONSTRAINT evm_forwarders_address_key UNIQUE (address);


--
-- Name: forwarders evm_forwarders_pkey; Type: CONSTRAINT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE ONLY evm.forwarders
    ADD CONSTRAINT evm_forwarders_pkey PRIMARY KEY (id);


--
-- Name: log_poller_filters evm_log_poller_filters_pkey; Type: CONSTRAINT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE ONLY evm.log_poller_filters
    ADD CONSTRAINT evm_log_poller_filters_pkey PRIMARY KEY (id);


--
-- Name: upkeep_states evm_upkeep_states_pkey; Type: CONSTRAINT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE ONLY evm.upkeep_states
    ADD CONSTRAINT evm_upkeep_states_pkey PRIMARY KEY (id);


--
-- Name: heads heads_pkey1; Type: CONSTRAINT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE ONLY evm.heads
    ADD CONSTRAINT heads_pkey1 PRIMARY KEY (id);


--
-- Name: log_poller_blocks log_poller_blocks_pkey; Type: CONSTRAINT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE ONLY evm.log_poller_blocks
    ADD CONSTRAINT log_poller_blocks_pkey PRIMARY KEY (block_number, evm_chain_id);


--
-- Name: logs logs_pkey; Type: CONSTRAINT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE ONLY evm.logs
    ADD CONSTRAINT logs_pkey PRIMARY KEY (block_hash, log_index, evm_chain_id);


--
-- Name: evm_logs_by_timestamp; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX evm_logs_by_timestamp ON evm.logs USING btree (evm_chain_id, address, event_sig, block_timestamp, block_number);


--
-- Name: evm_logs_idx; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX evm_logs_idx ON evm.logs USING btree (evm_chain_id, block_number, address, event_sig);


--
-- Name: evm_logs_idx_data_word_five; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX evm_logs_idx_data_word_five ON evm.logs USING btree (address, event_sig, evm_chain_id, "substring"(data, 129, 32));


--
-- Name: evm_logs_idx_data_word_four; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX evm_logs_idx_data_word_four ON evm.logs USING btree (SUBSTRING(data FROM 97 FOR 32));


--
-- Name: evm_logs_idx_data_word_one; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX evm_logs_idx_data_word_one ON evm.logs USING btree (SUBSTRING(data FROM 1 FOR 32));


--
-- Name: evm_logs_idx_data_word_three; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX evm_logs_idx_data_word_three ON evm.logs USING btree (SUBSTRING(data FROM 65 FOR 32));


--
-- Name: evm_logs_idx_data_word_two; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX evm_logs_idx_data_word_two ON evm.logs USING btree (SUBSTRING(data FROM 33 FOR 32));


--
-- Name: evm_logs_idx_topic_four; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX evm_logs_idx_topic_four ON evm.logs USING btree ((topics[4]));


--
-- Name: evm_logs_idx_topic_three; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX evm_logs_idx_topic_three ON evm.logs USING btree ((topics[3]));


--
-- Name: evm_logs_idx_topic_two; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX evm_logs_idx_topic_two ON evm.logs USING btree ((topics[2]));


--
-- Name: evm_logs_idx_tx_hash; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX evm_logs_idx_tx_hash ON evm.logs USING btree (tx_hash);


--
-- Name: idx_eth_receipts_block_number; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX idx_eth_receipts_block_number ON evm.receipts USING btree (block_number);


--
-- Name: idx_eth_receipts_created_at; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX idx_eth_receipts_created_at ON evm.receipts USING brin (created_at);


--
-- Name: idx_eth_receipts_unique; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE UNIQUE INDEX idx_eth_receipts_unique ON evm.receipts USING btree (tx_hash, block_hash);


--
-- Name: idx_eth_tx_attempts_broadcast_before_block_num; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX idx_eth_tx_attempts_broadcast_before_block_num ON evm.tx_attempts USING btree (broadcast_before_block_num);


--
-- Name: idx_eth_tx_attempts_created_at; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX idx_eth_tx_attempts_created_at ON evm.tx_attempts USING brin (created_at);


--
-- Name: idx_eth_tx_attempts_hash; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE UNIQUE INDEX idx_eth_tx_attempts_hash ON evm.tx_attempts USING btree (hash);


--
-- Name: idx_eth_tx_attempts_unbroadcast; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX idx_eth_tx_attempts_unbroadcast ON evm.tx_attempts USING btree (state) WHERE (state <> 'broadcast'::public.eth_tx_attempts_state);


--
-- Name: idx_eth_tx_attempts_unique_gas_prices; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE UNIQUE INDEX idx_eth_tx_attempts_unique_gas_prices ON evm.tx_attempts USING btree (eth_tx_id, gas_price);


--
-- Name: idx_eth_txes_broadcast_at; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX idx_eth_txes_broadcast_at ON evm.txes USING brin (broadcast_at);


--
-- Name: idx_eth_txes_created_at; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX idx_eth_txes_created_at ON evm.txes USING brin (created_at);


--
-- Name: idx_eth_txes_from_address; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX idx_eth_txes_from_address ON evm.txes USING btree (from_address);


--
-- Name: idx_eth_txes_initial_broadcast_at; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX idx_eth_txes_initial_broadcast_at ON evm.txes USING brin (initial_broadcast_at);


--
-- Name: idx_eth_txes_min_unconfirmed_nonce_for_key_evm_chain_id; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX idx_eth_txes_min_unconfirmed_nonce_for_key_evm_chain_id ON evm.txes USING btree (evm_chain_id, from_address, nonce) WHERE (state = 'unconfirmed'::public.eth_txes_state);


--
-- Name: idx_eth_txes_nonce_from_address_per_evm_chain_id; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE UNIQUE INDEX idx_eth_txes_nonce_from_address_per_evm_chain_id ON evm.txes USING btree (evm_chain_id, from_address, nonce);


--
-- Name: idx_eth_txes_pipeline_run_task_id; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE UNIQUE INDEX idx_eth_txes_pipeline_run_task_id ON evm.txes USING btree (pipeline_task_run_id) WHERE (pipeline_task_run_id IS NOT NULL);


--
-- Name: idx_eth_txes_state_from_address_evm_chain_id; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX idx_eth_txes_state_from_address_evm_chain_id ON evm.txes USING btree (evm_chain_id, from_address, state) WHERE (state <> 'confirmed'::public.eth_txes_state);


--
-- Name: idx_eth_txes_unstarted_subject_id_evm_chain_id; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX idx_eth_txes_unstarted_subject_id_evm_chain_id ON evm.txes USING btree (evm_chain_id, subject, id) WHERE ((subject IS NOT NULL) AND (state = 'unstarted'::public.eth_txes_state));


--
-- Name: idx_evm_key_states_address; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX idx_evm_key_states_address ON evm.key_states USING btree (address);


--
-- Name: idx_evm_key_states_evm_chain_id_address; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE UNIQUE INDEX idx_evm_key_states_evm_chain_id_address ON evm.key_states USING btree (evm_chain_id, address);


--
-- Name: idx_evm_log_poller_blocks_order_by_block; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX idx_evm_log_poller_blocks_order_by_block ON evm.log_poller_blocks USING btree (evm_chain_id, block_number DESC);


--
-- Name: idx_evm_logs_ordered_by_block_and_created_at; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX idx_evm_logs_ordered_by_block_and_created_at ON evm.logs USING btree (evm_chain_id, address, event_sig, block_number, created_at);


--
-- Name: idx_evm_upkeep_state_added_at_chain_id; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX idx_evm_upkeep_state_added_at_chain_id ON evm.upkeep_states USING btree (evm_chain_id, inserted_at);


--
-- Name: idx_evm_upkeep_state_chainid_workid; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE UNIQUE INDEX idx_evm_upkeep_state_chainid_workid ON evm.upkeep_states USING btree (evm_chain_id, work_id);


--
-- Name: idx_forwarders_created_at; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX idx_forwarders_created_at ON evm.forwarders USING brin (created_at);


--
-- Name: idx_forwarders_evm_address; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX idx_forwarders_evm_address ON evm.forwarders USING btree (address);


--
-- Name: idx_forwarders_evm_chain_id; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX idx_forwarders_evm_chain_id ON evm.forwarders USING btree (evm_chain_id);


--
-- Name: idx_forwarders_updated_at; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX idx_forwarders_updated_at ON evm.forwarders USING brin (updated_at);


--
-- Name: idx_heads_evm_chain_id_hash; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE UNIQUE INDEX idx_heads_evm_chain_id_hash ON evm.heads USING btree (evm_chain_id, hash);


--
-- Name: idx_heads_evm_chain_id_number; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE INDEX idx_heads_evm_chain_id_number ON evm.heads USING btree (evm_chain_id, number);


--
-- Name: idx_only_one_in_progress_tx_per_account_id_per_evm_chain_id; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE UNIQUE INDEX idx_only_one_in_progress_tx_per_account_id_per_evm_chain_id ON evm.txes USING btree (evm_chain_id, from_address) WHERE (state = 'in_progress'::public.eth_txes_state);


--
-- Name: idx_only_one_unbroadcast_attempt_per_eth_tx; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE UNIQUE INDEX idx_only_one_unbroadcast_attempt_per_eth_tx ON evm.tx_attempts USING btree (eth_tx_id) WHERE (state <> 'broadcast'::public.eth_tx_attempts_state);


--
-- Name: log_poller_filters_hash_key; Type: INDEX; Schema: evm; Owner: chainlink_dev
--

CREATE UNIQUE INDEX log_poller_filters_hash_key ON evm.log_poller_filters USING btree (evm.f_log_poller_filter_hash(name, evm_chain_id, address, event, topic2, topic3, topic4));


--
-- Name: receipts eth_receipts_tx_hash_fkey; Type: FK CONSTRAINT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE ONLY evm.receipts
    ADD CONSTRAINT eth_receipts_tx_hash_fkey FOREIGN KEY (tx_hash) REFERENCES evm.tx_attempts(hash) ON DELETE CASCADE;


--
-- Name: tx_attempts eth_tx_attempts_eth_tx_id_fkey; Type: FK CONSTRAINT; Schema: evm; Owner: chainlink_dev
--

ALTER TABLE ONLY evm.tx_attempts
    ADD CONSTRAINT eth_tx_attempts_eth_tx_id_fkey FOREIGN KEY (eth_tx_id) REFERENCES evm.txes(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

