
ALTER TABLE ONLY public.bridge_types
    ADD CONSTRAINT bridge_types_pkey PRIMARY KEY (name);

CREATE TABLE public.bridge_types (
    name text NOT NULL,
    url text NOT NULL,
    confirmations bigint DEFAULT 0 NOT NULL,
    incoming_token_hash text NOT NULL,
    salt text NOT NULL,
    outgoing_token text NOT NULL,
    minimum_contract_payment character varying(255),
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);

CREATE INDEX idx_bridge_types_created_at ON public.bridge_types USING brin (created_at);

CREATE INDEX idx_bridge_types_updated_at ON public.bridge_types USING brin (updated_at);

ALTER TABLE ONLY public.configurations
    ADD CONSTRAINT configurations_name_key UNIQUE (name);

ALTER TABLE ONLY public.configurations
    ADD CONSTRAINT configurations_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.configurations ALTER COLUMN id SET DEFAULT nextval('public.configurations_id_seq'::regclass);

CREATE TABLE public.configurations (
    id bigint NOT NULL,
    name text NOT NULL,
    value text NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone
);

ALTER SEQUENCE public.configurations_id_seq OWNED BY public.configurations.id;

CREATE SEQUENCE public.configurations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE INDEX idx_configurations_name ON public.configurations USING btree (name);

ALTER TABLE ONLY public.cron_specs
    ADD CONSTRAINT cron_specs_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.cron_specs ALTER COLUMN id SET DEFAULT nextval('public.cron_specs_id_seq'::regclass);

CREATE TABLE public.cron_specs (
    id integer NOT NULL,
    cron_schedule text NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);

CREATE SEQUENCE public.cron_specs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.cron_specs_id_seq OWNED BY public.cron_specs.id;

ALTER TABLE ONLY public.csa_keys
    ADD CONSTRAINT csa_keys_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.csa_keys
    ADD CONSTRAINT csa_keys_public_key_key UNIQUE (public_key);

ALTER TABLE ONLY public.csa_keys ALTER COLUMN id SET DEFAULT nextval('public.csa_keys_id_seq'::regclass);

CREATE TABLE public.csa_keys (
    id bigint NOT NULL,
    public_key bytea NOT NULL,
    encrypted_private_key jsonb NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    CONSTRAINT csa_keys_public_key_check CHECK ((octet_length(public_key) = 32))
);

ALTER SEQUENCE public.csa_keys_id_seq OWNED BY public.csa_keys.id;

CREATE SEQUENCE public.csa_keys_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE ONLY public.direct_request_specs
    ADD CONSTRAINT direct_request_specs_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES public.evm_chains(id) DEFERRABLE;

ALTER TABLE ONLY public.direct_request_specs
    ADD CONSTRAINT direct_request_specs_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.direct_request_specs ALTER COLUMN id SET DEFAULT nextval('public.eth_request_event_specs_id_seq'::regclass);

CREATE TABLE public.direct_request_specs (
    id integer NOT NULL,
    contract_address bytea NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    min_incoming_confirmations bigint,
    requesters text,
    min_contract_payment numeric(78,0),
    evm_chain_id numeric(78,0),
    CONSTRAINT eth_request_event_specs_contract_address_check CHECK ((octet_length(contract_address) = 20))
);

CREATE TABLE public.encrypted_key_rings (
    encrypted_keys jsonb,
    updated_at timestamp with time zone NOT NULL
);

ALTER TABLE ONLY public.encrypted_ocr_key_bundles
    ADD CONSTRAINT encrypted_ocr_key_bundles_pkey PRIMARY KEY (id);

CREATE TABLE public.encrypted_ocr_key_bundles (
    id bytea NOT NULL,
    on_chain_signing_address bytea NOT NULL,
    off_chain_public_key bytea NOT NULL,
    encrypted_private_keys jsonb NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    config_public_key bytea NOT NULL,
    deleted_at timestamp with time zone
);

ALTER TABLE ONLY public.encrypted_p2p_keys
    ADD CONSTRAINT encrypted_p2p_keys_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.encrypted_p2p_keys ALTER COLUMN id SET DEFAULT nextval('public.encrypted_p2p_keys_id_seq'::regclass);

CREATE TABLE public.encrypted_p2p_keys (
    id integer NOT NULL,
    peer_id text NOT NULL,
    pub_key bytea NOT NULL,
    encrypted_priv_key jsonb NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone,
    CONSTRAINT chk_pub_key_length CHECK ((octet_length(pub_key) = 32))
);

CREATE SEQUENCE public.encrypted_p2p_keys_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.encrypted_p2p_keys_id_seq OWNED BY public.encrypted_p2p_keys.id;

CREATE UNIQUE INDEX idx_unique_peer_ids ON public.encrypted_p2p_keys USING btree (peer_id);

CREATE UNIQUE INDEX idx_unique_pub_keys ON public.encrypted_p2p_keys USING btree (pub_key);

ALTER TABLE ONLY public.encrypted_vrf_keys
    ADD CONSTRAINT encrypted_secret_keys_pkey PRIMARY KEY (public_key);

CREATE TABLE public.encrypted_vrf_keys (
    public_key character varying(68) NOT NULL,
    vrf_key text NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone
);

ALTER TABLE ONLY public.eth_key_states
    ADD CONSTRAINT eth_key_states_address_key UNIQUE (address);

ALTER TABLE ONLY public.eth_key_states
    ADD CONSTRAINT eth_key_states_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES public.evm_chains(id) DEFERRABLE;

ALTER TABLE ONLY public.eth_key_states
    ADD CONSTRAINT eth_key_states_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.eth_key_states ALTER COLUMN id SET DEFAULT nextval('public.eth_key_states_id_seq'::regclass);

CREATE TABLE public.eth_key_states (
    id integer NOT NULL,
    address bytea NOT NULL,
    next_nonce bigint DEFAULT 0 NOT NULL,
    is_funding boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    evm_chain_id numeric(78,0) NOT NULL,
    CONSTRAINT chk_address_length CHECK ((octet_length(address) = 20))
);

CREATE SEQUENCE public.eth_key_states_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.eth_key_states_id_seq OWNED BY public.eth_key_states.id;

ALTER TABLE ONLY public.eth_receipts
    ADD CONSTRAINT eth_receipts_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.eth_receipts
    ADD CONSTRAINT eth_receipts_tx_hash_fkey FOREIGN KEY (tx_hash) REFERENCES public.eth_tx_attempts(hash) ON DELETE CASCADE;

ALTER TABLE ONLY public.eth_receipts ALTER COLUMN id SET DEFAULT nextval('public.eth_receipts_id_seq'::regclass);

CREATE TABLE public.eth_receipts (
    id bigint NOT NULL,
    tx_hash bytea NOT NULL,
    block_hash bytea NOT NULL,
    block_number bigint NOT NULL,
    transaction_index bigint NOT NULL,
    receipt jsonb NOT NULL,
    created_at timestamp with time zone NOT NULL,
    CONSTRAINT chk_hash_length CHECK (((octet_length(tx_hash) = 32) AND (octet_length(block_hash) = 32)))
);

ALTER SEQUENCE public.eth_receipts_id_seq OWNED BY public.eth_receipts.id;

CREATE SEQUENCE public.eth_receipts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE INDEX idx_eth_receipts_created_at ON public.eth_receipts USING brin (created_at);

CREATE INDEX idx_eth_receipts_block_number ON public.eth_receipts USING btree (block_number);

CREATE UNIQUE INDEX idx_eth_receipts_unique ON public.eth_receipts USING btree (tx_hash, block_hash);

CREATE SEQUENCE public.eth_request_event_specs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.eth_request_event_specs_id_seq OWNED BY public.direct_request_specs.id;

ALTER TABLE public.eth_tx_attempts
    ADD CONSTRAINT chk_chain_specific_gas_limit_not_zero CHECK ((chain_specific_gas_limit > 0)) NOT VALID;

ALTER TABLE ONLY public.eth_tx_attempts
    ADD CONSTRAINT eth_tx_attempts_eth_tx_id_fkey FOREIGN KEY (eth_tx_id) REFERENCES public.eth_txes(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.eth_tx_attempts
    ADD CONSTRAINT eth_tx_attempts_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.eth_tx_attempts ALTER COLUMN id SET DEFAULT nextval('public.eth_tx_attempts_id_seq'::regclass);

CREATE TABLE public.eth_tx_attempts (
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
    CONSTRAINT chk_cannot_broadcast_before_block_zero CHECK (((broadcast_before_block_num IS NULL) OR (broadcast_before_block_num > 0))),
    CONSTRAINT chk_eth_tx_attempts_fsm CHECK ((((state = ANY (ARRAY['in_progress'::public.eth_tx_attempts_state, 'insufficient_eth'::public.eth_tx_attempts_state])) AND (broadcast_before_block_num IS NULL)) OR (state = 'broadcast'::public.eth_tx_attempts_state))),
    CONSTRAINT chk_hash_length CHECK ((octet_length(hash) = 32)),
    CONSTRAINT chk_legacy_or_dynamic CHECK ((((tx_type = 0) AND (gas_price IS NOT NULL) AND (gas_tip_cap IS NULL) AND (gas_fee_cap IS NULL)) OR ((tx_type = 2) AND (gas_price IS NULL) AND (gas_tip_cap IS NOT NULL) AND (gas_fee_cap IS NOT NULL)))),
    CONSTRAINT chk_sanity_fee_cap_tip_cap CHECK (((gas_tip_cap IS NULL) OR (gas_fee_cap IS NULL) OR (gas_tip_cap <= gas_fee_cap))),
    CONSTRAINT chk_signed_raw_tx_present CHECK ((octet_length(signed_raw_tx) > 0)),
    CONSTRAINT chk_tx_type_is_byte CHECK (((tx_type >= 0) AND (tx_type <= 255)))
);

ALTER SEQUENCE public.eth_tx_attempts_id_seq OWNED BY public.eth_tx_attempts.id;

CREATE SEQUENCE public.eth_tx_attempts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TYPE public.eth_tx_attempts_state AS ENUM (
    'in_progress',
    'insufficient_eth',
    'broadcast'
);

CREATE INDEX idx_eth_tx_attempts_created_at ON public.eth_tx_attempts USING brin (created_at);

CREATE INDEX idx_eth_tx_attempts_broadcast_before_block_num ON public.eth_tx_attempts USING btree (broadcast_before_block_num);

CREATE UNIQUE INDEX idx_eth_tx_attempts_unique_gas_prices ON public.eth_tx_attempts USING btree (eth_tx_id, gas_price);

CREATE UNIQUE INDEX idx_only_one_unbroadcast_attempt_per_eth_tx ON public.eth_tx_attempts USING btree (eth_tx_id) WHERE (state <> 'broadcast'::public.eth_tx_attempts_state);

CREATE UNIQUE INDEX idx_eth_tx_attempts_hash ON public.eth_tx_attempts USING btree (hash);

CREATE INDEX idx_eth_tx_attempts_unbroadcast ON public.eth_tx_attempts USING btree (state) WHERE (state <> 'broadcast'::public.eth_tx_attempts_state);

ALTER TABLE ONLY public.eth_txes
    ADD CONSTRAINT eth_txes_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES public.evm_chains(id) DEFERRABLE;

ALTER TABLE ONLY public.eth_txes
    ADD CONSTRAINT eth_txes_from_address_fkey FOREIGN KEY (from_address) REFERENCES public.eth_key_states(address) NOT VALID;

ALTER TABLE ONLY public.eth_txes
    ADD CONSTRAINT eth_txes_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.eth_txes ALTER COLUMN id SET DEFAULT nextval('public.eth_txes_id_seq'::regclass);

CREATE TABLE public.eth_txes (
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
    access_list jsonb,
    simulate boolean DEFAULT false NOT NULL,
    CONSTRAINT chk_broadcast_at_is_sane CHECK ((broadcast_at > '2019-01-01 00:00:00+00'::timestamp with time zone)),
    CONSTRAINT chk_error_cannot_be_empty CHECK (((error IS NULL) OR (length(error) > 0))),
    CONSTRAINT chk_eth_txes_fsm CHECK ((((state = 'unstarted'::public.eth_txes_state) AND (nonce IS NULL) AND (error IS NULL) AND (broadcast_at IS NULL)) OR ((state = 'in_progress'::public.eth_txes_state) AND (nonce IS NOT NULL) AND (error IS NULL) AND (broadcast_at IS NULL)) OR ((state = 'fatal_error'::public.eth_txes_state) AND (nonce IS NULL) AND (error IS NOT NULL) AND (broadcast_at IS NULL)) OR ((state = 'unconfirmed'::public.eth_txes_state) AND (nonce IS NOT NULL) AND (error IS NULL) AND (broadcast_at IS NOT NULL)) OR ((state = 'confirmed'::public.eth_txes_state) AND (nonce IS NOT NULL) AND (error IS NULL) AND (broadcast_at IS NOT NULL)) OR ((state = 'confirmed_missing_receipt'::public.eth_txes_state) AND (nonce IS NOT NULL) AND (error IS NULL) AND (broadcast_at IS NOT NULL)))),
    CONSTRAINT chk_from_address_length CHECK ((octet_length(from_address) = 20)),
    CONSTRAINT chk_to_address_length CHECK ((octet_length(to_address) = 20))
);

ALTER SEQUENCE public.eth_txes_id_seq OWNED BY public.eth_txes.id;

CREATE SEQUENCE public.eth_txes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TYPE public.eth_txes_state AS ENUM (
    'unstarted',
    'in_progress',
    'fatal_error',
    'unconfirmed',
    'confirmed_missing_receipt',
    'confirmed'
);

CREATE INDEX idx_eth_txes_broadcast_at ON public.eth_txes USING brin (broadcast_at);

CREATE INDEX idx_eth_txes_created_at ON public.eth_txes USING brin (created_at);

CREATE UNIQUE INDEX idx_eth_txes_nonce_from_address_per_evm_chain_id ON public.eth_txes USING btree (evm_chain_id, from_address, nonce);

CREATE INDEX idx_eth_txes_min_unconfirmed_nonce_for_key_evm_chain_id ON public.eth_txes USING btree (evm_chain_id, from_address, nonce) WHERE (state = 'unconfirmed'::public.eth_txes_state);

CREATE INDEX idx_eth_txes_state_from_address_evm_chain_id ON public.eth_txes USING btree (evm_chain_id, from_address, state) WHERE (state <> 'confirmed'::public.eth_txes_state);

CREATE UNIQUE INDEX idx_only_one_in_progress_tx_per_account_id_per_evm_chain_id ON public.eth_txes USING btree (evm_chain_id, from_address) WHERE (state = 'in_progress'::public.eth_txes_state);

CREATE INDEX idx_eth_txes_unstarted_subject_id_evm_chain_id ON public.eth_txes USING btree (evm_chain_id, subject, id) WHERE ((subject IS NOT NULL) AND (state = 'unstarted'::public.eth_txes_state));

CREATE INDEX idx_eth_txes_from_address ON public.eth_txes USING btree (from_address);

CREATE UNIQUE INDEX idx_eth_txes_pipeline_run_task_id ON public.eth_txes USING btree (pipeline_task_run_id) WHERE (pipeline_task_run_id IS NOT NULL);

ALTER TABLE ONLY public.evm_chains
    ADD CONSTRAINT evm_chains_pkey PRIMARY KEY (id);

CREATE TABLE public.evm_chains (
    id numeric(78,0) NOT NULL,
    cfg jsonb DEFAULT '{}'::jsonb NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    enabled boolean DEFAULT true NOT NULL
);

ALTER TABLE ONLY public.external_initiators
    ADD CONSTRAINT access_key_unique UNIQUE (access_key);

ALTER TABLE ONLY public.external_initiators
    ADD CONSTRAINT external_initiators_name_unique UNIQUE (name);

ALTER TABLE ONLY public.external_initiators
    ADD CONSTRAINT external_initiators_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.external_initiators ALTER COLUMN id SET DEFAULT nextval('public.external_initiators_id_seq'::regclass);

CREATE TABLE public.external_initiators (
    id bigint NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    name text NOT NULL,
    url text,
    access_key text NOT NULL,
    salt text NOT NULL,
    hashed_secret text NOT NULL,
    outgoing_secret text NOT NULL,
    outgoing_token text NOT NULL
);

CREATE SEQUENCE public.external_initiators_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.external_initiators_id_seq OWNED BY public.external_initiators.id;

CREATE UNIQUE INDEX external_initiators_name_key ON public.external_initiators USING btree (lower(name));

ALTER TABLE ONLY public.external_initiator_webhook_specs
    ADD CONSTRAINT external_initiator_webhook_specs_external_initiator_id_fkey FOREIGN KEY (external_initiator_id) REFERENCES public.external_initiators(id) ON DELETE RESTRICT DEFERRABLE;

ALTER TABLE ONLY public.external_initiator_webhook_specs
    ADD CONSTRAINT external_initiator_webhook_specs_pkey PRIMARY KEY (external_initiator_id, webhook_spec_id);

ALTER TABLE ONLY public.external_initiator_webhook_specs
    ADD CONSTRAINT external_initiator_webhook_specs_webhook_spec_id_fkey FOREIGN KEY (webhook_spec_id) REFERENCES public.webhook_specs(id) ON DELETE CASCADE DEFERRABLE;

CREATE TABLE public.external_initiator_webhook_specs (
    external_initiator_id bigint NOT NULL,
    webhook_spec_id integer NOT NULL,
    spec jsonb NOT NULL
);

CREATE INDEX idx_external_initiator_webhook_specs_webhook_spec_id ON public.external_initiator_webhook_specs USING btree (webhook_spec_id);

ALTER TABLE ONLY public.feeds_managers
    ADD CONSTRAINT feeds_managers_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.feeds_managers
    ADD CONSTRAINT feeds_managers_public_key_key UNIQUE (public_key);

ALTER TABLE ONLY public.feeds_managers ALTER COLUMN id SET DEFAULT nextval('public.feeds_managers_id_seq'::regclass);

CREATE TABLE public.feeds_managers (
    id bigint NOT NULL,
    name character varying(255) NOT NULL,
    uri character varying(255) NOT NULL,
    public_key bytea NOT NULL,
    job_types text[] NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    is_ocr_bootstrap_peer boolean DEFAULT false NOT NULL,
    ocr_bootstrap_peer_multiaddr character varying,
    CONSTRAINT chk_ocr_bootstrap_peer_multiaddr CHECK ((NOT (is_ocr_bootstrap_peer AND ((ocr_bootstrap_peer_multiaddr IS NULL) OR ((ocr_bootstrap_peer_multiaddr)::text = ''::text))))),
    CONSTRAINT feeds_managers_public_key_check CHECK ((octet_length(public_key) = 32))
);

ALTER SEQUENCE public.feeds_managers_id_seq OWNED BY public.feeds_managers.id;

CREATE SEQUENCE public.feeds_managers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE ONLY public.flux_monitor_round_stats_v2
    ADD CONSTRAINT flux_monitor_round_stats_v2_aggregator_round_id_key UNIQUE (aggregator, round_id);

ALTER TABLE ONLY public.flux_monitor_round_stats_v2
    ADD CONSTRAINT flux_monitor_round_stats_v2_pipeline_run_id_fkey FOREIGN KEY (pipeline_run_id) REFERENCES public.pipeline_runs(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.flux_monitor_round_stats_v2
    ADD CONSTRAINT flux_monitor_round_stats_v2_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.flux_monitor_round_stats_v2 ALTER COLUMN id SET DEFAULT nextval('public.flux_monitor_round_stats_v2_id_seq'::regclass);

CREATE TABLE public.flux_monitor_round_stats_v2 (
    id bigint NOT NULL,
    aggregator bytea NOT NULL,
    round_id integer NOT NULL,
    num_new_round_logs integer DEFAULT 0 NOT NULL,
    num_submissions integer DEFAULT 0 NOT NULL,
    pipeline_run_id bigint
);

ALTER SEQUENCE public.flux_monitor_round_stats_v2_id_seq OWNED BY public.flux_monitor_round_stats_v2.id;

CREATE SEQUENCE public.flux_monitor_round_stats_v2_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE ONLY public.flux_monitor_specs
    ADD CONSTRAINT flux_monitor_specs_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES public.evm_chains(id) DEFERRABLE;

ALTER TABLE ONLY public.flux_monitor_specs
    ADD CONSTRAINT flux_monitor_specs_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.flux_monitor_specs ALTER COLUMN id SET DEFAULT nextval('public.flux_monitor_specs_id_seq'::regclass);

CREATE TABLE public.flux_monitor_specs (
    id integer NOT NULL,
    contract_address bytea NOT NULL,
    threshold real,
    absolute_threshold real,
    poll_timer_period bigint,
    poll_timer_disabled boolean,
    idle_timer_period bigint,
    idle_timer_disabled boolean,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    min_payment numeric(78,0),
    drumbeat_enabled boolean DEFAULT false NOT NULL,
    drumbeat_schedule text DEFAULT ''::text NOT NULL,
    drumbeat_random_delay bigint DEFAULT 0 NOT NULL,
    evm_chain_id numeric(78,0),
    CONSTRAINT flux_monitor_specs_check CHECK ((poll_timer_disabled OR (poll_timer_period > 0))),
    CONSTRAINT flux_monitor_specs_check1 CHECK ((idle_timer_disabled OR (idle_timer_period > 0))),
    CONSTRAINT flux_monitor_specs_contract_address_check CHECK ((octet_length(contract_address) = 20))
);

CREATE SEQUENCE public.flux_monitor_specs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.flux_monitor_specs_id_seq OWNED BY public.flux_monitor_specs.id;

ALTER TABLE ONLY public.goose_migrations
    ADD CONSTRAINT goose_migrations_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.goose_migrations ALTER COLUMN id SET DEFAULT nextval('public.goose_migrations_id_seq'::regclass);

CREATE TABLE public.goose_migrations (
    id integer NOT NULL,
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp without time zone DEFAULT now()
);

CREATE SEQUENCE public.goose_migrations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.goose_migrations_id_seq OWNED BY public.goose_migrations.id;

ALTER TABLE ONLY public.heads
    ADD CONSTRAINT heads_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES public.evm_chains(id) DEFERRABLE;

ALTER TABLE ONLY public.heads
    ADD CONSTRAINT heads_pkey1 PRIMARY KEY (id);

ALTER TABLE ONLY public.heads ALTER COLUMN id SET DEFAULT nextval('public.heads_id_seq'::regclass);

CREATE TABLE public.heads (
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

ALTER SEQUENCE public.heads_id_seq OWNED BY public.heads.id;

CREATE SEQUENCE public.heads_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE UNIQUE INDEX idx_heads_evm_chain_id_hash ON public.heads USING btree (evm_chain_id, hash);

CREATE INDEX idx_heads_evm_chain_id_number ON public.heads USING btree (evm_chain_id, number);

ALTER TABLE ONLY public.job_proposals
    ADD CONSTRAINT fk_feeds_manager FOREIGN KEY (feeds_manager_id) REFERENCES public.feeds_managers(id) DEFERRABLE;

ALTER TABLE ONLY public.job_proposals
    ADD CONSTRAINT job_proposals_job_id_fkey FOREIGN KEY (external_job_id) REFERENCES public.jobs(external_job_id) DEFERRABLE;

ALTER TABLE ONLY public.job_proposals
    ADD CONSTRAINT job_proposals_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.job_proposals ALTER COLUMN id SET DEFAULT nextval('public.job_proposals_id_seq'::regclass);

CREATE TABLE public.job_proposals (
    id bigint NOT NULL,
    spec text NOT NULL,
    status public.job_proposal_status NOT NULL,
    external_job_id uuid,
    feeds_manager_id integer NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    remote_uuid uuid NOT NULL,
    multiaddrs text[],
    proposed_at timestamp with time zone NOT NULL,
    CONSTRAINT chk_job_proposals_status_fsm CHECK ((((status = 'pending'::public.job_proposal_status) AND (external_job_id IS NULL)) OR ((status = 'approved'::public.job_proposal_status) AND (external_job_id IS NOT NULL)) OR ((status = 'rejected'::public.job_proposal_status) AND (external_job_id IS NULL)) OR ((status = 'cancelled'::public.job_proposal_status) AND (external_job_id IS NULL))))
);

ALTER SEQUENCE public.job_proposals_id_seq OWNED BY public.job_proposals.id;

CREATE SEQUENCE public.job_proposals_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TYPE public.job_proposal_status AS ENUM (
    'pending',
    'approved',
    'rejected',
    'cancelled'
);

CREATE UNIQUE INDEX idx_job_proposals_external_job_id ON public.job_proposals USING btree (external_job_id);

CREATE INDEX idx_job_proposals_feeds_manager_id ON public.job_proposals USING btree (feeds_manager_id);

CREATE UNIQUE INDEX idx_job_proposals_remote_uuid ON public.job_proposals USING btree (remote_uuid);

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT external_job_id_uniq UNIQUE (external_job_id);

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT jobs_cron_spec_id_fkey FOREIGN KEY (cron_spec_id) REFERENCES public.cron_specs(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT jobs_direct_request_spec_id_fkey FOREIGN KEY (direct_request_spec_id) REFERENCES public.direct_request_specs(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT jobs_flux_monitor_spec_id_fkey FOREIGN KEY (flux_monitor_spec_id) REFERENCES public.flux_monitor_specs(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT jobs_keeper_spec_id_fkey FOREIGN KEY (keeper_spec_id) REFERENCES public.keeper_specs(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT jobs_offchainreporting_oracle_spec_id_fkey FOREIGN KEY (offchainreporting_oracle_spec_id) REFERENCES public.offchainreporting_oracle_specs(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT jobs_pipeline_spec_id_fkey FOREIGN KEY (pipeline_spec_id) REFERENCES public.pipeline_specs(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT jobs_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT jobs_vrf_spec_id_fkey FOREIGN KEY (vrf_spec_id) REFERENCES public.vrf_specs(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT jobs_webhook_spec_id_fkey FOREIGN KEY (webhook_spec_id) REFERENCES public.webhook_specs(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.jobs ALTER COLUMN id SET DEFAULT nextval('public.jobs_id_seq'::regclass);

CREATE TABLE public.jobs (
    id integer NOT NULL,
    pipeline_spec_id integer NOT NULL,
    offchainreporting_oracle_spec_id integer,
    name character varying(255),
    schema_version integer NOT NULL,
    type character varying(255) NOT NULL,
    max_task_duration bigint,
    direct_request_spec_id integer,
    flux_monitor_spec_id integer,
    keeper_spec_id integer,
    cron_spec_id integer,
    vrf_spec_id integer,
    webhook_spec_id integer,
    external_job_id uuid NOT NULL,
    created_at timestamp with time zone NOT NULL,
    CONSTRAINT chk_only_one_spec CHECK ((num_nonnulls(offchainreporting_oracle_spec_id, direct_request_spec_id, flux_monitor_spec_id, keeper_spec_id, cron_spec_id, vrf_spec_id, webhook_spec_id) = 1)),
    CONSTRAINT chk_schema_version CHECK ((schema_version > 0)),
    CONSTRAINT chk_type CHECK (((type)::text <> ''::text)),
    CONSTRAINT non_zero_uuid_check CHECK ((external_job_id <> '00000000-0000-0000-0000-000000000000'::uuid))
);

CREATE SEQUENCE public.jobs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.jobs_id_seq OWNED BY public.jobs.id;

ALTER TABLE ONLY public.job_spec_errors
    ADD CONSTRAINT job_spec_errors_v2_job_id_fkey FOREIGN KEY (job_id) REFERENCES public.jobs(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.job_spec_errors
    ADD CONSTRAINT job_spec_errors_v2_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.job_spec_errors ALTER COLUMN id SET DEFAULT nextval('public.job_spec_errors_v2_id_seq'::regclass);

CREATE TABLE public.job_spec_errors (
    id bigint NOT NULL,
    job_id integer,
    description text NOT NULL,
    occurrences integer DEFAULT 1 NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);

CREATE INDEX idx_job_spec_errors_v2_created_at ON public.job_spec_errors USING brin (created_at);

CREATE INDEX idx_job_spec_errors_v2_finished_at ON public.job_spec_errors USING brin (updated_at);

CREATE UNIQUE INDEX job_spec_errors_v2_unique_idx ON public.job_spec_errors USING btree (job_id, description);

ALTER SEQUENCE public.job_spec_errors_v2_id_seq OWNED BY public.job_spec_errors.id;

CREATE SEQUENCE public.job_spec_errors_v2_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE INDEX idx_jobs_created_at ON public.jobs USING brin (created_at);

CREATE UNIQUE INDEX idx_jobs_unique_cron_spec_id ON public.jobs USING btree (cron_spec_id);

CREATE UNIQUE INDEX idx_jobs_unique_direct_request_spec_id ON public.jobs USING btree (direct_request_spec_id);

CREATE UNIQUE INDEX idx_jobs_unique_flux_monitor_spec_id ON public.jobs USING btree (flux_monitor_spec_id);

CREATE UNIQUE INDEX idx_jobs_unique_keeper_spec_id ON public.jobs USING btree (keeper_spec_id);

CREATE UNIQUE INDEX idx_jobs_name ON public.jobs USING btree (name);

CREATE UNIQUE INDEX idx_jobs_unique_offchain_reporting_oracle_spec_id ON public.jobs USING btree (offchainreporting_oracle_spec_id);

CREATE UNIQUE INDEX idx_jobs_unique_pipeline_spec_id ON public.jobs USING btree (pipeline_spec_id);

CREATE UNIQUE INDEX idx_jobs_unique_vrf_spec_id ON public.jobs USING btree (vrf_spec_id);

CREATE UNIQUE INDEX idx_jobs_unique_webhook_spec_id ON public.jobs USING btree (webhook_spec_id);

ALTER TABLE ONLY public.keeper_registries
    ADD CONSTRAINT keeper_registries_contract_address_key UNIQUE (contract_address);

ALTER TABLE ONLY public.keeper_registries
    ADD CONSTRAINT keeper_registries_job_id_fkey FOREIGN KEY (job_id) REFERENCES public.jobs(id) ON DELETE CASCADE DEFERRABLE;

ALTER TABLE ONLY public.keeper_registries
    ADD CONSTRAINT keeper_registries_job_id_key UNIQUE (job_id);

ALTER TABLE ONLY public.keeper_registries
    ADD CONSTRAINT keeper_registries_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.keeper_registries ALTER COLUMN id SET DEFAULT nextval('public.keeper_registries_id_seq'::regclass);

CREATE TABLE public.keeper_registries (
    id bigint NOT NULL,
    job_id integer NOT NULL,
    keeper_index integer NOT NULL,
    contract_address bytea NOT NULL,
    from_address bytea NOT NULL,
    check_gas integer NOT NULL,
    block_count_per_turn integer NOT NULL,
    num_keepers integer NOT NULL,
    CONSTRAINT keeper_registries_contract_address_check CHECK ((octet_length(contract_address) = 20)),
    CONSTRAINT keeper_registries_from_address_check CHECK ((octet_length(from_address) = 20))
);

ALTER SEQUENCE public.keeper_registries_id_seq OWNED BY public.keeper_registries.id;

CREATE SEQUENCE public.keeper_registries_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE INDEX idx_keeper_registries_keeper_index ON public.keeper_registries USING btree (keeper_index);

ALTER TABLE ONLY public.keeper_specs
    ADD CONSTRAINT keeper_specs_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES public.evm_chains(id) DEFERRABLE;

ALTER TABLE ONLY public.keeper_specs
    ADD CONSTRAINT keeper_specs_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.keeper_specs ALTER COLUMN id SET DEFAULT nextval('public.keeper_specs_id_seq'::regclass);

CREATE TABLE public.keeper_specs (
    id bigint NOT NULL,
    contract_address bytea NOT NULL,
    from_address bytea NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    evm_chain_id numeric(78,0),
    min_incoming_confirmations integer,
    CONSTRAINT keeper_specs_contract_address_check CHECK ((octet_length(contract_address) = 20)),
    CONSTRAINT keeper_specs_from_address_check CHECK ((octet_length(from_address) = 20))
);

ALTER SEQUENCE public.keeper_specs_id_seq OWNED BY public.keeper_specs.id;

CREATE SEQUENCE public.keeper_specs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE ONLY public.keys
    ADD CONSTRAINT keys_pkey PRIMARY KEY (id);

CREATE TABLE public.keys (
    address bytea NOT NULL,
    json jsonb NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    next_nonce bigint DEFAULT 0 NOT NULL,
    id integer NOT NULL,
    is_funding boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    CONSTRAINT chk_address_length CHECK ((octet_length(address) = 20))
);

ALTER TABLE ONLY public.keys ALTER COLUMN id SET DEFAULT nextval('public.keys_id_seq'::regclass);

CREATE SEQUENCE public.keys_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.keys_id_seq OWNED BY public.keys.id;

CREATE UNIQUE INDEX idx_unique_keys_address ON public.keys USING btree (address);

CREATE UNIQUE INDEX idx_keys_only_one_funding ON public.keys USING btree (is_funding) WHERE (is_funding = true);

CREATE TABLE public.lease_lock (
    client_id uuid NOT NULL,
    expires_at timestamp with time zone NOT NULL
);

CREATE UNIQUE INDEX only_one_lease_lock ON public.lease_lock USING btree (((client_id IS NOT NULL)));

ALTER TABLE ONLY public.log_broadcasts
    ADD CONSTRAINT log_broadcasts_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES public.evm_chains(id) DEFERRABLE;

ALTER TABLE ONLY public.log_broadcasts
    ADD CONSTRAINT log_consumptions_job_id_v2_fkey FOREIGN KEY (job_id) REFERENCES public.jobs(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.log_broadcasts
    ADD CONSTRAINT log_consumptions_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.log_broadcasts ALTER COLUMN id SET DEFAULT nextval('public.log_consumptions_id_seq'::regclass);

CREATE TABLE public.log_broadcasts (
    id bigint NOT NULL,
    block_hash bytea NOT NULL,
    log_index bigint NOT NULL,
    created_at timestamp without time zone NOT NULL,
    block_number bigint,
    job_id integer,
    consumed boolean DEFAULT false NOT NULL,
    evm_chain_id numeric(78,0) NOT NULL,
    updated_at timestamp with time zone NOT NULL
);

ALTER TABLE ONLY public.log_broadcasts_pending
    ADD CONSTRAINT log_broadcasts_pending_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES public.evm_chains(id) DEFERRABLE;

ALTER TABLE ONLY public.log_broadcasts_pending
    ADD CONSTRAINT log_broadcasts_pending_pkey PRIMARY KEY (evm_chain_id);

CREATE TABLE public.log_broadcasts_pending (
    evm_chain_id numeric(78,0) NOT NULL,
    block_number bigint,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);

CREATE INDEX log_consumptions_created_at_idx ON public.log_broadcasts USING brin (created_at);

CREATE INDEX idx_log_broadcasts_unconsumed ON public.log_broadcasts USING btree (evm_chain_id, block_number) WHERE ((consumed = false) AND (block_number IS NOT NULL));

CREATE UNIQUE INDEX log_broadcasts_unique_idx ON public.log_broadcasts USING btree (job_id, block_hash, log_index, evm_chain_id);

CREATE INDEX idx_log_broadcasts_unconsumed_job_id_v2 ON public.log_broadcasts USING btree (job_id, evm_chain_id) WHERE ((consumed = false) AND (job_id IS NOT NULL));

ALTER TABLE ONLY public.log_configs
    ADD CONSTRAINT log_configs_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.log_configs
    ADD CONSTRAINT log_configs_service_name_key UNIQUE (service_name);

ALTER TABLE ONLY public.log_configs ALTER COLUMN id SET DEFAULT nextval('public.log_configs_id_seq'::regclass);

CREATE TABLE public.log_configs (
    id bigint NOT NULL,
    service_name text NOT NULL,
    log_level public.log_level NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

ALTER SEQUENCE public.log_configs_id_seq OWNED BY public.log_configs.id;

CREATE SEQUENCE public.log_configs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.log_consumptions_id_seq OWNED BY public.log_broadcasts.id;

CREATE SEQUENCE public.log_consumptions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TYPE public.log_level AS ENUM (
    'debug',
    'info',
    'warn',
    'error',
    'panic'
);

ALTER TABLE ONLY public.nodes
    ADD CONSTRAINT nodes_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES public.evm_chains(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.nodes
    ADD CONSTRAINT nodes_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.nodes ALTER COLUMN id SET DEFAULT nextval('public.nodes_id_seq'::regclass);

CREATE TABLE public.nodes (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    evm_chain_id numeric(78,0) NOT NULL,
    ws_url text,
    http_url text,
    send_only boolean NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    CONSTRAINT nodes_http_url_check CHECK ((http_url <> ''::text)),
    CONSTRAINT nodes_name_check CHECK (((name)::text <> ''::text)),
    CONSTRAINT nodes_ws_url_check CHECK ((ws_url <> ''::text)),
    CONSTRAINT primary_or_sendonly CHECK (((send_only AND (ws_url IS NULL) AND (http_url IS NOT NULL)) OR ((NOT send_only) AND (ws_url IS NOT NULL))))
);

CREATE SEQUENCE public.nodes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.nodes_id_seq OWNED BY public.nodes.id;

CREATE INDEX idx_nodes_evm_chain_id ON public.nodes USING btree (evm_chain_id);

CREATE UNIQUE INDEX idx_nodes_unique_name ON public.nodes USING btree (lower((name)::text));

ALTER TABLE ONLY public.node_versions
    ADD CONSTRAINT node_versions_pkey PRIMARY KEY (version);

CREATE UNIQUE INDEX idx_only_one_node_version ON public.node_versions USING btree (((version IS NOT NULL)));

CREATE TABLE public.node_versions (
    version text NOT NULL,
    created_at timestamp without time zone NOT NULL
);

ALTER TABLE ONLY public.offchainreporting_contract_configs
    ADD CONSTRAINT offchainreporting_contract_configs_pkey PRIMARY KEY (offchainreporting_oracle_spec_id);

ALTER TABLE ONLY public.offchainreporting_contract_configs
    ADD CONSTRAINT offchainreporting_contract_co_offchainreporting_oracle_spe_fkey FOREIGN KEY (offchainreporting_oracle_spec_id) REFERENCES public.offchainreporting_oracle_specs(id) ON DELETE CASCADE;

CREATE TABLE public.offchainreporting_contract_configs (
    offchainreporting_oracle_spec_id integer NOT NULL,
    config_digest bytea NOT NULL,
    signers bytea[],
    transmitters bytea[],
    threshold integer,
    encoded_config_version bigint,
    encoded bytea,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    CONSTRAINT offchainreporting_contract_configs_config_digest_check CHECK ((octet_length(config_digest) = 16))
);

ALTER TABLE ONLY public.offchainreporting_discoverer_announcements
    ADD CONSTRAINT offchainreporting_discoverer_announcements_local_peer_id_fkey FOREIGN KEY (local_peer_id) REFERENCES public.encrypted_p2p_keys(peer_id) DEFERRABLE;

ALTER TABLE ONLY public.offchainreporting_discoverer_announcements
    ADD CONSTRAINT offchainreporting_discoverer_announcements_pkey PRIMARY KEY (local_peer_id, remote_peer_id);

CREATE TABLE public.offchainreporting_discoverer_announcements (
    local_peer_id text NOT NULL,
    remote_peer_id text NOT NULL,
    ann bytea NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);

ALTER TABLE ONLY public.offchainreporting_latest_round_requested
    ADD CONSTRAINT offchainreporting_latest_round_requested_pkey PRIMARY KEY (offchainreporting_oracle_spec_id);

ALTER TABLE ONLY public.offchainreporting_latest_round_requested
    ADD CONSTRAINT offchainreporting_latest_roun_offchainreporting_oracle_spe_fkey FOREIGN KEY (offchainreporting_oracle_spec_id) REFERENCES public.offchainreporting_oracle_specs(id) ON DELETE CASCADE DEFERRABLE;

CREATE TABLE public.offchainreporting_latest_round_requested (
    offchainreporting_oracle_spec_id integer NOT NULL,
    requester bytea NOT NULL,
    config_digest bytea NOT NULL,
    epoch bigint NOT NULL,
    round bigint NOT NULL,
    raw jsonb NOT NULL,
    CONSTRAINT offchainreporting_latest_round_requested_config_digest_check CHECK ((octet_length(config_digest) = 16)),
    CONSTRAINT offchainreporting_latest_round_requested_requester_check CHECK ((octet_length(requester) = 20))
);

ALTER TABLE ONLY public.offchainreporting_oracle_specs
    ADD CONSTRAINT offchainreporting_oracle_specs_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES public.evm_chains(id) DEFERRABLE;

ALTER TABLE ONLY public.offchainreporting_oracle_specs
    ADD CONSTRAINT offchainreporting_oracle_specs_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.offchainreporting_oracle_specs ALTER COLUMN id SET DEFAULT nextval('public.offchainreporting_oracle_specs_id_seq'::regclass);

CREATE TABLE public.offchainreporting_oracle_specs (
    id integer NOT NULL,
    contract_address bytea NOT NULL,
    p2p_peer_id text,
    p2p_bootstrap_peers text[],
    is_bootstrap_peer boolean NOT NULL,
    encrypted_ocr_key_bundle_id bytea,
    transmitter_address bytea,
    observation_timeout bigint,
    blockchain_timeout bigint,
    contract_config_tracker_subscribe_interval bigint,
    contract_config_tracker_poll_interval bigint,
    contract_config_confirmations integer NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    evm_chain_id numeric(78,0),
    CONSTRAINT chk_contract_address_length CHECK ((octet_length(contract_address) = 20))
);

CREATE SEQUENCE public.offchainreporting_oracle_specs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.offchainreporting_oracle_specs_id_seq OWNED BY public.offchainreporting_oracle_specs.id;

CREATE INDEX idx_offchainreporting_oracle_specs_created_at ON public.offchainreporting_oracle_specs USING brin (created_at);

CREATE INDEX idx_offchainreporting_oracle_specs_updated_at ON public.offchainreporting_oracle_specs USING brin (updated_at);

CREATE UNIQUE INDEX unique_contract_addr_per_chain ON public.offchainreporting_oracle_specs USING btree (contract_address, evm_chain_id) WHERE (evm_chain_id IS NOT NULL);

CREATE UNIQUE INDEX unique_contract_addr ON public.offchainreporting_oracle_specs USING btree (contract_address) WHERE (evm_chain_id IS NULL);

ALTER TABLE ONLY public.offchainreporting_pending_transmissions
    ADD CONSTRAINT offchainreporting_pending_transmissions_pkey PRIMARY KEY (offchainreporting_oracle_spec_id, config_digest, epoch, round);

ALTER TABLE ONLY public.offchainreporting_pending_transmissions
    ADD CONSTRAINT offchainreporting_pending_tra_offchainreporting_oracle_spe_fkey FOREIGN KEY (offchainreporting_oracle_spec_id) REFERENCES public.offchainreporting_oracle_specs(id) ON DELETE CASCADE;

CREATE TABLE public.offchainreporting_pending_transmissions (
    offchainreporting_oracle_spec_id integer NOT NULL,
    config_digest bytea NOT NULL,
    epoch bigint NOT NULL,
    round bigint NOT NULL,
    "time" timestamp with time zone NOT NULL,
    median numeric(78,0) NOT NULL,
    serialized_report bytea NOT NULL,
    rs bytea[] NOT NULL,
    ss bytea[] NOT NULL,
    vs bytea NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    CONSTRAINT offchainreporting_pending_transmissions_config_digest_check CHECK ((octet_length(config_digest) = 16))
);

CREATE INDEX idx_offchainreporting_pending_transmissions_time ON public.offchainreporting_pending_transmissions USING btree ("time");

ALTER TABLE ONLY public.offchainreporting_persistent_states
    ADD CONSTRAINT offchainreporting_persistent__offchainreporting_oracle_spe_fkey FOREIGN KEY (offchainreporting_oracle_spec_id) REFERENCES public.offchainreporting_oracle_specs(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.offchainreporting_persistent_states
    ADD CONSTRAINT offchainreporting_persistent_states_pkey PRIMARY KEY (offchainreporting_oracle_spec_id, config_digest);

CREATE TABLE public.offchainreporting_persistent_states (
    offchainreporting_oracle_spec_id integer NOT NULL,
    config_digest bytea NOT NULL,
    epoch bigint NOT NULL,
    highest_sent_epoch bigint NOT NULL,
    highest_received_epoch bigint[] NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    CONSTRAINT offchainreporting_persistent_states_config_digest_check CHECK ((octet_length(config_digest) = 16))
);

CREATE TABLE public.p2p_peers (
    id text NOT NULL,
    addr text NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    peer_id text NOT NULL
);

CREATE INDEX p2p_peers_id ON public.p2p_peers USING btree (id);

CREATE INDEX p2p_peers_peer_id ON public.p2p_peers USING btree (peer_id);

ALTER TABLE ONLY public.pipeline_runs
    ADD CONSTRAINT pipeline_runs_pipeline_spec_id_fkey FOREIGN KEY (pipeline_spec_id) REFERENCES public.pipeline_specs(id) ON DELETE CASCADE DEFERRABLE;

ALTER TABLE ONLY public.pipeline_runs
    ADD CONSTRAINT pipeline_runs_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.pipeline_runs ALTER COLUMN id SET DEFAULT nextval('public.pipeline_runs_id_seq'::regclass);

CREATE TABLE public.pipeline_runs (
    id bigint NOT NULL,
    pipeline_spec_id integer NOT NULL,
    meta jsonb DEFAULT '{}'::jsonb,
    created_at timestamp with time zone NOT NULL,
    finished_at timestamp with time zone,
    fatal_errors jsonb,
    outputs jsonb,
    inputs jsonb,
    state public.pipeline_runs_state DEFAULT 'completed'::public.pipeline_runs_state NOT NULL,
    all_errors jsonb,
    CONSTRAINT pipeline_runs_check CHECK ((((state = ANY (ARRAY['completed'::public.pipeline_runs_state, 'errored'::public.pipeline_runs_state])) AND (finished_at IS NOT NULL) AND (num_nulls(outputs, fatal_errors) = 0)) OR ((state = ANY (ARRAY['running'::public.pipeline_runs_state, 'suspended'::public.pipeline_runs_state])) AND (num_nulls(finished_at, outputs, fatal_errors) = 3))))
);

ALTER SEQUENCE public.pipeline_runs_id_seq OWNED BY public.pipeline_runs.id;

CREATE SEQUENCE public.pipeline_runs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TYPE public.pipeline_runs_state AS ENUM (
    'running',
    'suspended',
    'errored',
    'completed'
);

CREATE INDEX idx_pipeline_runs_created_at ON public.pipeline_runs USING brin (created_at);

CREATE INDEX idx_pipeline_runs_finished_at ON public.pipeline_runs USING brin (finished_at);

CREATE INDEX idx_pipeline_runs_unfinished_runs ON public.pipeline_runs USING btree (id) WHERE (finished_at IS NULL);

CREATE INDEX pipeline_runs_suspended ON public.pipeline_runs USING btree (id) WHERE (state = 'suspended'::public.pipeline_runs_state);

CREATE INDEX idx_pipeline_runs_pipeline_spec_id ON public.pipeline_runs USING btree (pipeline_spec_id);

ALTER TABLE ONLY public.pipeline_specs
    ADD CONSTRAINT pipeline_specs_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.pipeline_specs ALTER COLUMN id SET DEFAULT nextval('public.pipeline_specs_id_seq'::regclass);

CREATE TABLE public.pipeline_specs (
    id integer NOT NULL,
    dot_dag_source text NOT NULL,
    created_at timestamp with time zone NOT NULL,
    max_task_duration bigint
);

CREATE SEQUENCE public.pipeline_specs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.pipeline_specs_id_seq OWNED BY public.pipeline_specs.id;

CREATE INDEX idx_pipeline_specs_created_at ON public.pipeline_specs USING brin (created_at);

ALTER TABLE ONLY public.pipeline_task_runs
    ADD CONSTRAINT pipeline_task_runs_pipeline_run_id_fkey FOREIGN KEY (pipeline_run_id) REFERENCES public.pipeline_runs(id) ON DELETE CASCADE DEFERRABLE;

ALTER TABLE ONLY public.pipeline_task_runs
    ADD CONSTRAINT pipeline_task_runs_pkey PRIMARY KEY (id);

CREATE TABLE public.pipeline_task_runs (
    pipeline_run_id bigint NOT NULL,
    type text NOT NULL,
    index integer DEFAULT 0 NOT NULL,
    output jsonb,
    error text,
    created_at timestamp with time zone NOT NULL,
    finished_at timestamp with time zone,
    dot_id text NOT NULL,
    id uuid NOT NULL,
    CONSTRAINT chk_pipeline_task_run_fsm CHECK ((((finished_at IS NOT NULL) AND (num_nonnulls(output, error) <> 2)) OR (num_nulls(finished_at, output, error) = 3)))
);

CREATE INDEX idx_pipeline_task_runs_created_at ON public.pipeline_task_runs USING brin (created_at);

CREATE INDEX idx_pipeline_task_runs_finished_at ON public.pipeline_task_runs USING brin (finished_at);

CREATE UNIQUE INDEX pipeline_task_runs_pipeline_run_id_dot_id_idx ON public.pipeline_task_runs USING btree (pipeline_run_id, dot_id);

CREATE INDEX idx_unfinished_pipeline_task_runs ON public.pipeline_task_runs USING btree (pipeline_run_id) WHERE (finished_at IS NULL);

CREATE TYPE public.run_status AS ENUM (
    'unstarted',
    'in_progress',
    'pending_incoming_confirmations',
    'pending_outgoing_confirmations',
    'pending_connection',
    'pending_bridge',
    'pending_sleep',
    'errored',
    'completed',
    'cancelled'
);

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);

CREATE TABLE public.sessions (
    id text NOT NULL,
    last_used timestamp with time zone,
    created_at timestamp with time zone NOT NULL
);

CREATE INDEX idx_sessions_created_at ON public.sessions USING brin (created_at);

CREATE INDEX idx_sessions_last_used ON public.sessions USING brin (last_used);

ALTER TABLE ONLY public.upkeep_registrations
    ADD CONSTRAINT upkeep_registrations_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.upkeep_registrations
    ADD CONSTRAINT upkeep_registrations_registry_id_fkey FOREIGN KEY (registry_id) REFERENCES public.keeper_registries(id) ON DELETE CASCADE DEFERRABLE;

ALTER TABLE ONLY public.upkeep_registrations ALTER COLUMN id SET DEFAULT nextval('public.upkeep_registrations_id_seq'::regclass);

CREATE TABLE public.upkeep_registrations (
    id bigint NOT NULL,
    registry_id bigint NOT NULL,
    execute_gas integer NOT NULL,
    check_data bytea NOT NULL,
    upkeep_id bigint NOT NULL,
    positioning_constant integer NOT NULL,
    last_run_block_height bigint DEFAULT 0 NOT NULL
);

ALTER SEQUENCE public.upkeep_registrations_id_seq OWNED BY public.upkeep_registrations.id;

CREATE SEQUENCE public.upkeep_registrations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE UNIQUE INDEX idx_upkeep_registrations_unique_upkeep_ids_per_keeper ON public.upkeep_registrations USING btree (registry_id, upkeep_id);

CREATE INDEX idx_upkeep_registrations_upkeep_id ON public.upkeep_registrations USING btree (upkeep_id);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (email);

CREATE TABLE public.users (
    email text NOT NULL,
    hashed_password text,
    created_at timestamp with time zone NOT NULL,
    token_key text,
    token_salt text,
    token_hashed_secret text,
    updated_at timestamp with time zone NOT NULL
);

CREATE INDEX idx_users_updated_at ON public.users USING brin (updated_at);

CREATE INDEX idx_users_created_at ON public.users USING btree (created_at);

ALTER TABLE ONLY public.vrf_specs
    ADD CONSTRAINT vrf_specs_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES public.evm_chains(id) DEFERRABLE;

ALTER TABLE ONLY public.vrf_specs
    ADD CONSTRAINT vrf_specs_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.vrf_specs ALTER COLUMN id SET DEFAULT nextval('public.vrf_specs_id_seq'::regclass);

CREATE TABLE public.vrf_specs (
    id bigint NOT NULL,
    public_key text NOT NULL,
    coordinator_address bytea NOT NULL,
    min_incoming_confirmations bigint NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    evm_chain_id numeric(78,0),
    from_address bytea,
    poll_period bigint,
    CONSTRAINT coordinator_address_len_chk CHECK ((octet_length(coordinator_address) = 20))
);

ALTER SEQUENCE public.vrf_specs_id_seq OWNED BY public.vrf_specs.id;

CREATE SEQUENCE public.vrf_specs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE ONLY public.web_authns
    ADD CONSTRAINT fk_email FOREIGN KEY (email) REFERENCES public.users(email);

ALTER TABLE ONLY public.web_authns
    ADD CONSTRAINT web_authns_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.web_authns ALTER COLUMN id SET DEFAULT nextval('public.web_authns_id_seq'::regclass);

CREATE TABLE public.web_authns (
    id bigint NOT NULL,
    email text NOT NULL,
    public_key_data jsonb NOT NULL
);

ALTER SEQUENCE public.web_authns_id_seq OWNED BY public.web_authns.id;

CREATE SEQUENCE public.web_authns_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE UNIQUE INDEX web_authns_email_idx ON public.web_authns USING btree (lower(email));

ALTER TABLE ONLY public.webhook_specs
    ADD CONSTRAINT webhook_specs_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.webhook_specs ALTER COLUMN id SET DEFAULT nextval('public.webhook_specs_id_seq'::regclass);

CREATE TABLE public.webhook_specs (
    id integer NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);

CREATE SEQUENCE public.webhook_specs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.webhook_specs_id_seq OWNED BY public.webhook_specs.id;
