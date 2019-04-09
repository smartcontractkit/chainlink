--
-- PostgreSQL database dump
--

-- Dumped from database version 11.1
-- Dumped by pg_dump version 11.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: bridge_types; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bridge_types (
    name character varying(255) NOT NULL,
    url character varying(255),
    confirmations bigint,
    incoming_token character varying(255),
    outgoing_token character varying(255),
    minimum_contract_payment character varying(255)
);


--
-- Name: encumbrances; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.encumbrances (
    id integer NOT NULL,
    payment character varying(255),
    expiration bigint,
    end_at timestamp with time zone,
    oracles text
);


--
-- Name: encumbrances_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.encumbrances_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: encumbrances_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.encumbrances_id_seq OWNED BY public.encumbrances.id;


--
-- Name: heads; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.heads (
    hash character varying NOT NULL,
    number bigint NOT NULL
);


--
-- Name: initiators; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.initiators (
    id integer NOT NULL,
    job_spec_id character varying(36),
    type text NOT NULL,
    created_at timestamp with time zone,
    schedule text,
    "time" timestamp with time zone,
    ran boolean,
    address bytea,
    requesters text,
    deleted_at timestamp with time zone
);


--
-- Name: initiators_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.initiators_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: initiators_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.initiators_id_seq OWNED BY public.initiators.id;


--
-- Name: job_runs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.job_runs (
    id text NOT NULL,
    job_spec_id character varying(36) NOT NULL,
    result_id integer,
    run_request_id integer,
    status text,
    created_at timestamp with time zone,
    completed_at timestamp with time zone,
    updated_at timestamp with time zone,
    initiator_id integer,
    creation_height character varying(255),
    observed_height character varying(255),
    overrides_id integer,
    deleted_at timestamp with time zone
);


--
-- Name: job_specs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.job_specs (
    id text NOT NULL,
    created_at timestamp with time zone,
    start_at timestamp with time zone,
    end_at timestamp with time zone,
    deleted_at timestamp with time zone
);


--
-- Name: keys; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.keys (
    address character varying(64) NOT NULL,
    json text
);


--
-- Name: migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.migrations (
    id character varying(12) NOT NULL
);


--
-- Name: run_requests; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.run_requests (
    id integer NOT NULL,
    request_id text,
    tx_hash bytea,
    requester bytea,
    created_at timestamp with time zone
);


--
-- Name: run_requests_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.run_requests_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: run_requests_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.run_requests_id_seq OWNED BY public.run_requests.id;


--
-- Name: run_results; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.run_results (
    id integer NOT NULL,
    cached_job_run_id text,
    cached_task_run_id text,
    data text,
    status text,
    error_message text,
    amount character varying(255)
);


--
-- Name: run_results_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.run_results_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: run_results_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.run_results_id_seq OWNED BY public.run_results.id;


--
-- Name: service_agreements; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.service_agreements (
    id text NOT NULL,
    created_at timestamp with time zone,
    encumbrance_id integer,
    request_body text,
    signature character varying(255),
    job_spec_id character varying(36) NOT NULL
);


--
-- Name: sessions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sessions (
    id text NOT NULL,
    last_used timestamp with time zone,
    created_at timestamp with time zone
);


--
-- Name: sync_events; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sync_events (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    body text
);


--
-- Name: sync_events_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.sync_events_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: sync_events_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.sync_events_id_seq OWNED BY public.sync_events.id;


--
-- Name: task_runs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.task_runs (
    id text NOT NULL,
    job_run_id character varying(36) NOT NULL,
    result_id integer,
    status text,
    task_spec_id integer,
    minimum_confirmations bigint,
    created_at timestamp with time zone
);


--
-- Name: task_specs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.task_specs (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    job_spec_id character varying(36),
    type text NOT NULL,
    confirmations bigint,
    params text
);


--
-- Name: task_specs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.task_specs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: task_specs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.task_specs_id_seq OWNED BY public.task_specs.id;


--
-- Name: tx_attempts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tx_attempts (
    hash bytea NOT NULL,
    tx_id bigint,
    gas_price character varying(255),
    confirmed boolean,
    hex text,
    sent_at bigint,
    created_at timestamp with time zone
);


--
-- Name: txes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.txes (
    id bigint NOT NULL,
    "from" bytea NOT NULL,
    "to" bytea NOT NULL,
    data bytea,
    nonce bigint,
    value character varying(255),
    gas_limit bigint,
    hash bytea,
    gas_price character varying(255),
    confirmed boolean,
    hex text,
    sent_at bigint
);


--
-- Name: txes_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.txes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: txes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.txes_id_seq OWNED BY public.txes.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    email text NOT NULL,
    hashed_password text,
    created_at timestamp with time zone
);


--
-- Name: encumbrances id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.encumbrances ALTER COLUMN id SET DEFAULT nextval('public.encumbrances_id_seq'::regclass);


--
-- Name: initiators id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.initiators ALTER COLUMN id SET DEFAULT nextval('public.initiators_id_seq'::regclass);


--
-- Name: run_requests id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.run_requests ALTER COLUMN id SET DEFAULT nextval('public.run_requests_id_seq'::regclass);


--
-- Name: run_results id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.run_results ALTER COLUMN id SET DEFAULT nextval('public.run_results_id_seq'::regclass);


--
-- Name: sync_events id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sync_events ALTER COLUMN id SET DEFAULT nextval('public.sync_events_id_seq'::regclass);


--
-- Name: task_specs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.task_specs ALTER COLUMN id SET DEFAULT nextval('public.task_specs_id_seq'::regclass);


--
-- Name: txes id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.txes ALTER COLUMN id SET DEFAULT nextval('public.txes_id_seq'::regclass);


--
-- Data for Name: bridge_types; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.bridge_types (name, url, confirmations, incoming_token, outgoing_token, minimum_contract_payment) FROM stdin;
\.


--
-- Data for Name: encumbrances; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.encumbrances (id, payment, expiration, end_at, oracles) FROM stdin;
\.


--
-- Data for Name: heads; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.heads (hash, number) FROM stdin;
2639889a465c75b75d747a8c64b597ba58b9b1569c36390719ead14393e1caa6	124365
fb2f6931161a54f2fac09f7b68594b893614b147b2b5b319c0a84d0f4f314baf	124366
32ef05b89f508476a908a3f834452e5523d10e5cc608f3209e2057b0fcbda19a	124367
ee507e0c6c7abfeeb2711af4c171ef425e928b37501fd44586c8a727a454b88d	124368
fb11bb20dffcb1e3121563c0eb9e25d03ef47466254d4eb9b2c067435a9270c8	124369
0d6d569b25b7ee943b995b5feaf7f80baabf5b13be1ae5f181b492070f26b128	124370
9f963c5813cca84f68637e1b03d15d453554aa93b5ed51f55d34a22897e6761b	124371
3a869128de0b58ba1af83555905e5909fc9c2dfbcc5757e52097316f52928fdc	124372
cc3d503c1bb033ca678eb2f349df796a59896b3a6c89fd8abf86d71638cc6712	124373
70f055ed04de621077166019b8456c7ca023f2dc0fefe2f35b6835471b9e5eb9	124374
f22b426cac7d929a80ddb04d0c4f46125e3d819a69b5da77652477e630d35b42	124375
c2f8bcf4409e0efe5b0d8ea7ab2d4984738348c7b105c37071bb52e427e3a9af	124376
83f778d44e336a0116b45e6cb38547013cb42ba7a4c9879ae2414a92ae6e126f	124377
33472e1ed18a8f4b4a81521934e4fcd73f1447360c779010d6f4fbee8081c4e2	124378
d67c041ed041fe187894060b6c46e643e0f7d06b1e297f55f252ecafcfd04b18	124379
ebd8542c77361897d0be267647d7a31f9e436d3571bd973a0a40858c8156ff01	124380
9d552b13346d9e124805281957d7abc80ee8fd60f4537acf2be8dd85bed65748	124381
ef80b92847a44f1805db6acc0930e75f36154c1523f66a2f1c6136542128d86f	124382
19463c7d180ed79801331809d550715052da039313f8c842cf0f3ea189433a7b	124383
40cf82a21bf2ac72cbde914c5e7afe97975f9f0506d78b641a5aac617ae08921	124384
42ceba6c42a0fb76aa4c3f9be713073278cc148db1511e26789b162eed4daf2c	124385
61bf44ad12cc10f2518baa3f49c48c1fcaf2550ee81bc63e948cb19678ebb5cd	124386
29b62c15c3b072e3f0ab5c63737e9023f8a35defbb1f7a8d0e5166886b1a1f7c	124387
e26118c48249d2a2fc36f14a69b3b15c406401dda30d7de06e31175dc4a35875	124388
de52131a5548934de406de5b4bdadcb2561d8b54cc1f2803383919fe266a6209	124389
300ac3b765f88a4b173c2206e558154bd8048c839621e0476d27779289a049b0	124390
58ebecc469490be97baaffd0884ba31683cbc656a4ac6882e9258263b3bd163a	124391
bfc3e52c0a5bbcad6148cbfdc4942c6960874faecbaf84a7175798a54f74705a	124392
ca7753e16aa186059d94165642d6518edf7754bbc5a8cf3dc8e564d5426fdd6e	124393
47d1a7c6184405f95caf618087bd2063326acb0402661754c5985a93b3f56110	124394
d8c66b7bdbeb4ca0d85345227b38f4a32f807ba02024ef3359bc73fc508ccdf8	124395
a94d77e653205a724d2d9122de83c31ebd4fa8fd77b3c7e57ef94ad0b9e63ac2	124396
312183af37e91340b21571812a84694e1e6f8a39779b8e049ed11af18fca2ff9	124397
e3d063b9764825435f5a480f4b3b4d09d3f05eb551c6e7fdce35cc3b8f285516	124398
5512db0d9268e276a85c0f9e9d61824e8a5ccc71b1fa7a3b3e0063b630ff0703	124399
3b2fb5fe139ef83556819e85f04664f17923ebac935f86625001bfb8dc7fded6	124400
9e011b3af724b926af20079790607737c31e1b16742419916c91ed74220c61e3	124401
392fa199f7afbe95330d8eee606ba5972c17ce4ba37741bb40a800bb04540aab	124402
11aabaacbd9a2382823f574ea37c57156dd06c39bbafe44378ff776ada3343a1	124403
65114bc1f142d2f442b955894b9d7dcc33884e78d962c537fc18e2c99462729f	124404
4a06491560a4637974267dcabb49f2db5fb3efb6db6829a2dfb17798ace1ed6b	124405
28422d3ea4285958c8b8b74b421c78ceb5b1cda99df04e08e95361959b8717f1	124406
329bc867fddbf45adafa9b5879f7f4d38d2fc935d6e9410bbef94887e017a0f9	124407
cc5252ae01c9ebd42aa20247c5054fbe8cf01d4d15cefd73798b7615c134e31a	124408
534a40639a3a3ada62aa9ec637a3a2bb690d4f71add8b39c00113d3deea1042a	124409
6c2306559534f06e42103616fde45be5757be79fa613160b4e02cca654aa5ae3	124410
c25452fefce82d16a74e53a21387343cc8f43ca713814e7d5dd139cd37b947c6	124411
980788e52ca9cc676e66c3233662299ec3b92d1b0af2a0fc0be271b5c6de978f	124412
e8e4fe1696cfa80c33493ba1140d383f15bc975a6d752c1cec446593a66948b6	124413
310ec7c52077ca2aa1ead327e707236d64ca19ed245984e5fb097af75359eba4	124414
93c60b5848116f903011a219c831ae6032cf7e742c13f70e271d6f4e0bd9c030	124415
3d8d67f66c45ccdfef1028779e998a0b7b301867ab600984cc7358814edcc31c	124416
c626b77c7dc84689cfbcc90d09cce69dfe5269cd8dc9a018063d7facce8da385	124417
e3d3b82afd37541e5bbc1f97974af49493738345229e310c39984f9fcd4f6b1d	124418
019d0e9bac1fdc9680d4615e06b67113d17b824383b61b8f7af7fd57b7029ec6	124419
6a55dd953c30bd8b19a584e708ee0de44f0d601764072332f27c38c58f7756b8	124420
60f96fc8c423b98c9c867359c3305a2297035d40bd57ede9f7e2ff8eaf3f9f4f	124421
09bb31b6237a009eac7025a80cafc3622f1e27db42344559f59003bafe12615c	124422
06336176a564ff96f7447d597e489b5798231016f00e5439b409d5a7239facce	124423
bb011ef02747ba1af094660317ae1fc3a6fd27c2ed662dd77da1c5b8736918fd	124424
83f12a20f4ddc1e9b23e2410fa5c4fc3a37993c9424c30710dae0b75aa770bf1	124425
969ad07267d7d412b5f6d2cbaaa5c787ea7d1da899d0c893c517388b352cbd5c	124426
b81b16ddd867a73e901659b304b311b199a7c949835cbc3603768fca027b935a	124427
d88adf00a7b8c3b77ee8b203efedf563ec51999679776e75bb15618f57dee9b2	124428
d85591f0453809092c75d5e5997798caeb16cd7b195e579271818be2f5ca4d96	124429
936c0e306c3c92ee30f9db8794020a3f97ac3a6ae20871b8725f175e593bb19f	124430
311cd99c71aae9ca53ac782523afef564e0e80fa4223a9f8a3c5c1a83fe4b6f3	124431
c78777741907c1b6bc35dd7530494087e0860f717df6f8f4469711e7d1e24c37	124432
6e3cd539949866442515d1332f2ea2ee16b38828b157e555353813ed43d2c526	124433
75f5e02b371ee915a71b11d760ced6c01edc6178c2b0aca533b0ccc1d7d361c2	124434
a6bd928efa878305b6e2a7865c8eb85c7342545ddfdd94d406c336ab2818ee12	124435
f1cf4da843299ec42be3f128a1da9cc720d64578b7133c7d972be8ee37fba303	124436
5c62d3b6bdea56b3674523b63f5f3823bfd37aa81edb123bde0c648c31affaa7	124437
b88fdc1904ac2d41064372aff49de686eb6765bc6fec2e73d415dba8cd44c4d9	124438
984cec3e582d417bcf9ff325939bcf6103e27606c0a218298df325eb9bc73775	124439
63314299f0c3c5564e09b51ffae7b9d1704608338b86f3619cf0418c0677fdf7	124440
8ee9df64d96253f4236c59f45b3a882aa06d466fcbd77f602b80a46a5535d63b	124441
4ea3cbbd2aa3d3d48e9a0f987274fb98d94d2bdd91fe9cc97f87ec45a11dc4e9	124442
b64b4bb2b59e1d77f62e193388afc8acd3f3cf498deb06a1b45c6b6eea24cf89	124443
684c9d073dbb309837c61da46d4478140033723fcc1c14ee165bf5b8727692f1	124444
ee05fffbbcb23d433ab1ad42306e3ab2103ba1e6601fc38643faa195683e9316	124445
497d7b0fe30d475b4d8fd58414519197180fbcbae6c05f46a22018d2cfa59b71	124446
\.


--
-- Data for Name: initiators; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.initiators (id, job_spec_id, type, created_at, schedule, "time", ran, address, requesters, deleted_at) FROM stdin;
1	7b6e4281dba042c28a38758a34ed06ed	runlog	2019-04-08 14:11:17.521224-04		\N	f	\\x0000000000000000000000000000000000000000	0xEc4f6443e71c5131Ab5033e8fB46C6bCaF8b1b2B	\N
2	fe3793164b664937a371fff12c5fac01	ethlog	2019-04-08 14:11:41.121886-04		\N	f	\\x7d0e877e7fdd362a8c5249244fbbd437a1fc3032		\N
3	f439c385b52448ccafae205e7f49f37a	web	2019-04-08 14:11:53.429544-04		\N	f	\\x0000000000000000000000000000000000000000		\N
\.


--
-- Data for Name: job_runs; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.job_runs (id, job_spec_id, result_id, run_request_id, status, created_at, completed_at, updated_at, initiator_id, creation_height, observed_height, overrides_id, deleted_at) FROM stdin;
33cf659338ac4af7aa431c993cd30210	7b6e4281dba042c28a38758a34ed06ed	5	1	completed	2019-04-08 14:11:21.004589-04	2019-04-08 14:11:33.023895-04	2019-04-08 14:11:33.024868-04	1	124382	124388	2	\N
29cf7ba429414c95b954f6484757aa44	fe3793164b664937a371fff12c5fac01	8	2	completed	2019-04-08 14:11:45.005121-04	2019-04-08 14:11:45.013889-04	2019-04-08 14:11:45.014862-04	2	124394	124394	7	\N
7bf368af7d3e4ee9bf7146bce7150efb	f439c385b52448ccafae205e7f49f37a	13	0	completed	2019-04-08 14:11:53.545652-04	2019-04-08 14:11:59.020282-04	2019-04-08 14:11:59.02098-04	3	\N	124401	10	\N
\.


--
-- Data for Name: job_specs; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.job_specs (id, created_at, start_at, end_at, deleted_at) FROM stdin;
7b6e4281dba042c28a38758a34ed06ed	2019-04-08 14:11:17.52048-04	\N	\N	\N
fe3793164b664937a371fff12c5fac01	2019-04-08 14:11:41.121636-04	\N	\N	\N
f439c385b52448ccafae205e7f49f37a	2019-04-08 14:11:53.429206-04	\N	\N	\N
\.


--
-- Data for Name: keys; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.keys (address, json) FROM stdin;
0xc1e6ED650D11BBcB0b8bE72c548A0370D1F7ac79	{"address":"c1e6ed650d11bbcb0b8be72c548a0370d1f7ac79","crypto":{"cipher":"aes-128-ctr","ciphertext":"9732e18fc9b0ccd9b772bcb2632e658fa950583e9e1db4bb86ee98b0491e8a55","cipherparams":{"iv":"83ee401f5252c6fdb31e9e3f87fdfaea"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"a454ca06c8d216e015cf27a2981bd84bb7d703792365814349ac3d631fba5e5e"},"mac":"6b6748476b11e876f028cd9354b70cd116271d6c1fc305710a52e802490de01b"},"id":"af25bf97-5c99-4475-8700-8bc453c26b46","version":3}
\.


--
-- Data for Name: migrations; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.migrations (id) FROM stdin;
0
1549496047
1551816486
1551895034
1552418531
1553029703
1554131520
\.


--
-- Data for Name: run_requests; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.run_requests (id, request_id, tx_hash, requester, created_at) FROM stdin;
1	0xb35b1199fdf3a4f9ee87cb71c8cd7e7857ae395d93cd19cfa9e7f704aa3fb1cd	\\x9178a8fb2c92222ac6dc925bce71e76ea77a813031235846ea8923ed4ec9e1cb	\\xec4f6443e71c5131ab5033e8fb46c6bcaf8b1b2b	2019-04-08 14:11:21.015264-04
2	\N	\\xe7308e92f3168f6d9eef8316180f2dff85c5331adf27646415933028819fd371	\N	2019-04-08 14:11:45.005741-04
\.


--
-- Data for Name: run_results; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.run_results (id, cached_job_run_id, cached_task_run_id, data, status, error_message, amount) FROM stdin;
1			{"address":"0xa8fFA679D1f78D30928461D64e4c4bE92E8bDd35","dataPrefix":"0xb35b1199fdf3a4f9ee87cb71c8cd7e7857ae395d93cd19cfa9e7f704aa3fb1cd0000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000ec4f6443e71c5131ab5033e8fb46c6bcaf8b1b2b042f2b6500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005cab8ff5","functionSelector":"0x4ab0d190","msg":"hello_chainlink"}		\N	1000000000000000000
2			{"address":"0xa8fFA679D1f78D30928461D64e4c4bE92E8bDd35","dataPrefix":"0xb35b1199fdf3a4f9ee87cb71c8cd7e7857ae395d93cd19cfa9e7f704aa3fb1cd0000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000ec4f6443e71c5131ab5033e8fb46c6bcaf8b1b2b042f2b6500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005cab8ff5","functionSelector":"0x4ab0d190","msg":"hello_chainlink"}		\N	1000000000000000000
3	33cf659338ac4af7aa431c993cd30210	b193bdd969c94bc38ac66524acff140d	{"address":"0xa8fFA679D1f78D30928461D64e4c4bE92E8bDd35","dataPrefix":"0xb35b1199fdf3a4f9ee87cb71c8cd7e7857ae395d93cd19cfa9e7f704aa3fb1cd0000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000ec4f6443e71c5131ab5033e8fb46c6bcaf8b1b2b042f2b6500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005cab8ff5","functionSelector":"0x4ab0d190","msg":"hello_chainlink"}	pending_sleep	\N	\N
5	33cf659338ac4af7aa431c993cd30210	0f7ef39aa5c34fa1aed7f9579a3d4a92	{"address":"0xa8fFA679D1f78D30928461D64e4c4bE92E8bDd35","dataPrefix":"0xb35b1199fdf3a4f9ee87cb71c8cd7e7857ae395d93cd19cfa9e7f704aa3fb1cd0000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000ec4f6443e71c5131ab5033e8fb46c6bcaf8b1b2b042f2b6500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005cab8ff5","ethereumReceipts":[{"blockNumber":124387,"transactionHash":"0x0e7c006fb7d97cc78f3f9c6171fad8ac6acf443d9dc213f41a5c6757f26ec01d"}],"functionSelector":"0x4ab0d190","msg":"hello_chainlink","result":"0x0e7c006fb7d97cc78f3f9c6171fad8ac6acf443d9dc213f41a5c6757f26ec01d"}	completed	\N	\N
6			{"address":"0x7d0e877e7fdd362a8c5249244fbbd437a1fc3032","topics":["0xd226ed2bc8a4081ee6d62540525cd9a44aa022eca68e344a7f7b924a75faeed5","0x68656c6c6f5f636861696e6c696e6b0000000000000000000000000000000000"],"data":"0x","blockNumber":"0x1e5ea","transactionHash":"0xe7308e92f3168f6d9eef8316180f2dff85c5331adf27646415933028819fd371","transactionIndex":"0x0","blockHash":"0x47d1a7c6184405f95caf618087bd2063326acb0402661754c5985a93b3f56110","logIndex":"0x0","removed":false}		\N	\N
7			{"address":"0x7d0e877e7fdd362a8c5249244fbbd437a1fc3032","topics":["0xd226ed2bc8a4081ee6d62540525cd9a44aa022eca68e344a7f7b924a75faeed5","0x68656c6c6f5f636861696e6c696e6b0000000000000000000000000000000000"],"data":"0x","blockNumber":"0x1e5ea","transactionHash":"0xe7308e92f3168f6d9eef8316180f2dff85c5331adf27646415933028819fd371","transactionIndex":"0x0","blockHash":"0x47d1a7c6184405f95caf618087bd2063326acb0402661754c5985a93b3f56110","logIndex":"0x0","removed":false}		\N	\N
8	29cf7ba429414c95b954f6484757aa44	cdc439fb67d847648fe04d820e3cf9e2	{"address":"0x7d0e877e7fdd362a8c5249244fbbd437a1fc3032","blockHash":"0x47d1a7c6184405f95caf618087bd2063326acb0402661754c5985a93b3f56110","blockNumber":"0x1e5ea","data":"0x","logIndex":"0x0","removed":false,"result":"{\\"headers\\":{\\"host\\":\\"127.0.0.1:6690\\",\\"user-agent\\":\\"Go-http-client/1.1\\",\\"content-length\\":\\"467\\",\\"content-type\\":\\"application/json\\"},\\"body\\":{\\"address\\":\\"0x7d0e877e7fdd362a8c5249244fbbd437a1fc3032\\",\\"blockHash\\":\\"0x47d1a7c6184405f95caf618087bd2063326acb0402661754c5985a93b3f56110\\",\\"blockNumber\\":\\"0x1e5ea\\",\\"data\\":\\"0x\\",\\"logIndex\\":\\"0x0\\",\\"removed\\":false,\\"topics\\":[\\"0xd226ed2bc8a4081ee6d62540525cd9a44aa022eca68e344a7f7b924a75faeed5\\",\\"0x68656c6c6f5f636861696e6c696e6b0000000000000000000000000000000000\\"],\\"transactionHash\\":\\"0xe7308e92f3168f6d9eef8316180f2dff85c5331adf27646415933028819fd371\\",\\"transactionIndex\\":\\"0x0\\"}}","topics":["0xd226ed2bc8a4081ee6d62540525cd9a44aa022eca68e344a7f7b924a75faeed5","0x68656c6c6f5f636861696e6c696e6b0000000000000000000000000000000000"],"transactionHash":"0xe7308e92f3168f6d9eef8316180f2dff85c5331adf27646415933028819fd371","transactionIndex":"0x0"}	completed	\N	\N
4	33cf659338ac4af7aa431c993cd30210	b1e770326cf44891944a8f8dabf4b626	{"address":"0xa8fFA679D1f78D30928461D64e4c4bE92E8bDd35","dataPrefix":"0xb35b1199fdf3a4f9ee87cb71c8cd7e7857ae395d93cd19cfa9e7f704aa3fb1cd0000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000ec4f6443e71c5131ab5033e8fb46c6bcaf8b1b2b042f2b6500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005cab8ff5","functionSelector":"0x4ab0d190","msg":"hello_chainlink","result":"{\\"headers\\":{\\"host\\":\\"127.0.0.1:6690\\",\\"user-agent\\":\\"Go-http-client/1.1\\",\\"content-length\\":\\"450\\",\\"content-type\\":\\"application/json\\"},\\"body\\":{\\"address\\":\\"0xa8fFA679D1f78D30928461D64e4c4bE92E8bDd35\\",\\"dataPrefix\\":\\"0xb35b1199fdf3a4f9ee87cb71c8cd7e7857ae395d93cd19cfa9e7f704aa3fb1cd0000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000ec4f6443e71c5131ab5033e8fb46c6bcaf8b1b2b042f2b6500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005cab8ff5\\",\\"functionSelector\\":\\"0x4ab0d190\\",\\"msg\\":\\"hello_chainlink\\"}}"}	completed	\N	\N
9			{}		\N	\N
10			{}		\N	\N
11	7bf368af7d3e4ee9bf7146bce7150efb	de76106638224af0afb1b34a6dfb623f	{"result":"{\\"last\\": \\"3843.95\\"}"}	completed	\N	\N
12	7bf368af7d3e4ee9bf7146bce7150efb	167181e65aff4bd4939b0d2ad68fb066	{"result":"3843.95"}	completed	\N	\N
13	7bf368af7d3e4ee9bf7146bce7150efb	54ce916133c441b596b8821a497bca64	{"ethereumReceipts":[{"blockNumber":124400,"transactionHash":"0x272820b76613a5eacef327a0505b7319b0e04024419e90b217d96ae66ac93c72"}],"result":"0x272820b76613a5eacef327a0505b7319b0e04024419e90b217d96ae66ac93c72"}	completed	\N	\N
\.


--
-- Data for Name: service_agreements; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.service_agreements (id, created_at, encumbrance_id, request_body, signature, job_spec_id) FROM stdin;
\.


--
-- Data for Name: sessions; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.sessions (id, last_used, created_at) FROM stdin;
4f7a229e5f694fc58371c8ffe31b841c	2019-04-08 14:11:17.519207-04	2019-04-08 14:11:17.511328-04
e3acc2c32ea54819b7cf922c7c696cc4	2019-04-08 14:11:41.12049-04	2019-04-08 14:11:41.111321-04
f1db7ff333c34c02a390e538651b7004	2019-04-08 14:11:45.427575-04	2019-04-08 14:11:06.352503-04
2b55fcb9582d462abdf2611ea8b3563f	2019-04-08 14:11:53.670918-04	2019-04-08 14:11:50.695922-04
\.


--
-- Data for Name: sync_events; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.sync_events (id, created_at, updated_at, body) FROM stdin;
\.


--
-- Data for Name: task_runs; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.task_runs (id, job_run_id, result_id, status, task_spec_id, minimum_confirmations, created_at) FROM stdin;
de76106638224af0afb1b34a6dfb623f	7bf368af7d3e4ee9bf7146bce7150efb	11	completed	5	0	2019-04-08 14:11:53.548242-04
167181e65aff4bd4939b0d2ad68fb066	7bf368af7d3e4ee9bf7146bce7150efb	12	completed	6	0	2019-04-08 14:11:53.550249-04
54ce916133c441b596b8821a497bca64	7bf368af7d3e4ee9bf7146bce7150efb	13	completed	7	0	2019-04-08 14:11:53.55209-04
b193bdd969c94bc38ac66524acff140d	33cf659338ac4af7aa431c993cd30210	3	completed	1	0	2019-04-08 14:11:21.020935-04
b1e770326cf44891944a8f8dabf4b626	33cf659338ac4af7aa431c993cd30210	4	completed	2	0	2019-04-08 14:11:21.023587-04
0f7ef39aa5c34fa1aed7f9579a3d4a92	33cf659338ac4af7aa431c993cd30210	5	completed	3	0	2019-04-08 14:11:21.025124-04
cdc439fb67d847648fe04d820e3cf9e2	29cf7ba429414c95b954f6484757aa44	8	completed	4	0	2019-04-08 14:11:45.007899-04
\.


--
-- Data for Name: task_specs; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.task_specs (id, created_at, updated_at, deleted_at, job_spec_id, type, confirmations, params) FROM stdin;
1	2019-04-08 14:11:17.522027-04	2019-04-08 14:11:17.522027-04	\N	7b6e4281dba042c28a38758a34ed06ed	sleep	0	{"until":1554747087}
2	2019-04-08 14:11:17.522784-04	2019-04-08 14:11:17.522784-04	\N	7b6e4281dba042c28a38758a34ed06ed	httppost	0	{"url":"http://127.0.0.1:6690/count"}
3	2019-04-08 14:11:17.523087-04	2019-04-08 14:11:17.523087-04	\N	7b6e4281dba042c28a38758a34ed06ed	ethtx	0	{"functionSelector":"fulfillOracleRequest(uint256,bytes32)"}
4	2019-04-08 14:11:41.122187-04	2019-04-08 14:11:41.122187-04	\N	fe3793164b664937a371fff12c5fac01	httppost	0	{"url":"http://127.0.0.1:6690/count"}
5	2019-04-08 14:11:53.42991-04	2019-04-08 14:11:53.42991-04	\N	f439c385b52448ccafae205e7f49f37a	httpget	0	{"get":"http://localhost:55362"}
6	2019-04-08 14:11:53.430221-04	2019-04-08 14:11:53.430221-04	\N	f439c385b52448ccafae205e7f49f37a	jsonparse	0	{"path":["last"]}
7	2019-04-08 14:11:53.430557-04	2019-04-08 14:11:53.430557-04	\N	f439c385b52448ccafae205e7f49f37a	ethtx	0	{"address":"0xaa664fa2fdc390c662de1dbacf1218ac6e066ae6","functionSelector":"setBytes(bytes32,bytes)"}
\.


--
-- Data for Name: tx_attempts; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.tx_attempts (hash, tx_id, gas_price, confirmed, hex, sent_at, created_at) FROM stdin;
\\x0e7c006fb7d97cc78f3f9c6171fad8ac6acf443d9dc213f41a5c6757f26ec01d	1	20000000000	t	0xf9012c808504a817c8008307a12094a8ffa679d1f78d30928461d64e4c4be92e8bdd3580b8c44ab0d190b35b1199fdf3a4f9ee87cb71c8cd7e7857ae395d93cd19cfa9e7f704aa3fb1cd0000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000ec4f6443e71c5131ab5033e8fb46c6bcaf8b1b2b042f2b6500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005cab8ff50000000000000000000000000000000000000000000000000000000000000000820a96a0f99f220704644beb6ea89dc6a4a500319b466e22b5ea2024b10984bc45824fe2a0766948d3c4071b3c82b8c4c86adc07eec6cd02938022cbeb94e69292c21eb58d	124385	2019-04-08 14:11:27.038954-04
\\x272820b76613a5eacef327a0505b7319b0e04024419e90b217d96ae66ac93c72	2	20000000000	t	0xf88b018504a817c8008307a12094aa664fa2fdc390c662de1dbacf1218ac6e066ae680a42e28d0840000000000000000000000000000000000000000000000000000000000000384820a95a005515856466c4593d6bedf806bdde5927113d71dff45ecea79fccbdca667a02da02ea620c566c0b031df93d0e0bf3f12b1a78352e472bbdc3ac619559f456cd6fc	124398	2019-04-08 14:11:53.573541-04
\.


--
-- Data for Name: txes; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.txes (id, "from", "to", data, nonce, value, gas_limit, hash, gas_price, confirmed, hex, sent_at) FROM stdin;
1	\\xc1e6ed650d11bbcb0b8be72c548a0370d1f7ac79	\\xa8ffa679d1f78d30928461d64e4c4be92e8bdd35	\\x4ab0d190b35b1199fdf3a4f9ee87cb71c8cd7e7857ae395d93cd19cfa9e7f704aa3fb1cd0000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000ec4f6443e71c5131ab5033e8fb46c6bcaf8b1b2b042f2b6500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005cab8ff50000000000000000000000000000000000000000000000000000000000000000	0	0	500000	\\x0e7c006fb7d97cc78f3f9c6171fad8ac6acf443d9dc213f41a5c6757f26ec01d	20000000000	t	0xf9012c808504a817c8008307a12094a8ffa679d1f78d30928461d64e4c4be92e8bdd3580b8c44ab0d190b35b1199fdf3a4f9ee87cb71c8cd7e7857ae395d93cd19cfa9e7f704aa3fb1cd0000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000ec4f6443e71c5131ab5033e8fb46c6bcaf8b1b2b042f2b6500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005cab8ff50000000000000000000000000000000000000000000000000000000000000000820a96a0f99f220704644beb6ea89dc6a4a500319b466e22b5ea2024b10984bc45824fe2a0766948d3c4071b3c82b8c4c86adc07eec6cd02938022cbeb94e69292c21eb58d	124385
2	\\xc1e6ed650d11bbcb0b8be72c548a0370d1f7ac79	\\xaa664fa2fdc390c662de1dbacf1218ac6e066ae6	\\x2e28d0840000000000000000000000000000000000000000000000000000000000000384	1	0	500000	\\x272820b76613a5eacef327a0505b7319b0e04024419e90b217d96ae66ac93c72	20000000000	t	0xf88b018504a817c8008307a12094aa664fa2fdc390c662de1dbacf1218ac6e066ae680a42e28d0840000000000000000000000000000000000000000000000000000000000000384820a95a005515856466c4593d6bedf806bdde5927113d71dff45ecea79fccbdca667a02da02ea620c566c0b031df93d0e0bf3f12b1a78352e472bbdc3ac619559f456cd6fc	124398
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.users (email, hashed_password, created_at) FROM stdin;
notreal@fakeemail.ch	$2a$10$qVUatuT8/dVsieoIgH/vmux9nT2nhUWxujEwfq1AwRGMIOPr5XdJW	2019-04-08 14:10:45.04635-04
\.


--
-- Name: encumbrances_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.encumbrances_id_seq', 1, false);


--
-- Name: initiators_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.initiators_id_seq', 3, true);


--
-- Name: run_requests_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.run_requests_id_seq', 2, true);


--
-- Name: run_results_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.run_results_id_seq', 13, true);


--
-- Name: sync_events_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.sync_events_id_seq', 1, false);


--
-- Name: task_specs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.task_specs_id_seq', 7, true);


--
-- Name: txes_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.txes_id_seq', 2, true);


--
-- Name: bridge_types bridge_types_with_pk_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bridge_types
    ADD CONSTRAINT bridge_types_with_pk_pkey PRIMARY KEY (name);


--
-- Name: encumbrances encumbrances_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.encumbrances
    ADD CONSTRAINT encumbrances_pkey PRIMARY KEY (id);


--
-- Name: heads heads_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.heads
    ADD CONSTRAINT heads_pkey PRIMARY KEY (hash);


--
-- Name: initiators initiators_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.initiators
    ADD CONSTRAINT initiators_pkey PRIMARY KEY (id);


--
-- Name: job_runs job_runs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.job_runs
    ADD CONSTRAINT job_runs_pkey PRIMARY KEY (id);


--
-- Name: job_specs job_specs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.job_specs
    ADD CONSTRAINT job_specs_pkey PRIMARY KEY (id);


--
-- Name: keys keys_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.keys
    ADD CONSTRAINT keys_pkey PRIMARY KEY (address);


--
-- Name: migrations migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.migrations
    ADD CONSTRAINT migrations_pkey PRIMARY KEY (id);


--
-- Name: run_requests run_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.run_requests
    ADD CONSTRAINT run_requests_pkey PRIMARY KEY (id);


--
-- Name: run_results run_results_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.run_results
    ADD CONSTRAINT run_results_pkey PRIMARY KEY (id);


--
-- Name: service_agreements service_agreements_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.service_agreements
    ADD CONSTRAINT service_agreements_pkey PRIMARY KEY (id);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


--
-- Name: sync_events sync_events_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sync_events
    ADD CONSTRAINT sync_events_pkey PRIMARY KEY (id);


--
-- Name: task_runs task_runs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.task_runs
    ADD CONSTRAINT task_runs_pkey PRIMARY KEY (id);


--
-- Name: task_specs task_specs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.task_specs
    ADD CONSTRAINT task_specs_pkey PRIMARY KEY (id);


--
-- Name: tx_attempts tx_attempts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tx_attempts
    ADD CONSTRAINT tx_attempts_pkey PRIMARY KEY (hash);


--
-- Name: txes txes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.txes
    ADD CONSTRAINT txes_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (email);


--
-- Name: idx_heads_number; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_heads_number ON public.heads USING btree (number);


--
-- Name: idx_initiators_address; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_initiators_address ON public.initiators USING btree (address);


--
-- Name: idx_initiators_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_initiators_created_at ON public.initiators USING btree (created_at);


--
-- Name: idx_initiators_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_initiators_deleted_at ON public.initiators USING btree (deleted_at);


--
-- Name: idx_initiators_job_spec_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_initiators_job_spec_id ON public.initiators USING btree (job_spec_id);


--
-- Name: idx_initiators_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_initiators_type ON public.initiators USING btree (type);


--
-- Name: idx_job_runs_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_job_runs_created_at ON public.job_runs USING btree (created_at);


--
-- Name: idx_job_runs_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_job_runs_deleted_at ON public.job_runs USING btree (deleted_at);


--
-- Name: idx_job_runs_job_spec_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_job_runs_job_spec_id ON public.job_runs USING btree (job_spec_id);


--
-- Name: idx_job_runs_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_job_runs_status ON public.job_runs USING btree (status);


--
-- Name: idx_job_specs_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_job_specs_created_at ON public.job_specs USING btree (created_at);


--
-- Name: idx_job_specs_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_job_specs_deleted_at ON public.job_specs USING btree (deleted_at);


--
-- Name: idx_job_specs_end_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_job_specs_end_at ON public.job_specs USING btree (end_at);


--
-- Name: idx_job_specs_start_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_job_specs_start_at ON public.job_specs USING btree (start_at);


--
-- Name: idx_service_agreements_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_service_agreements_created_at ON public.service_agreements USING btree (created_at);


--
-- Name: idx_service_agreements_job_spec_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_service_agreements_job_spec_id ON public.service_agreements USING btree (job_spec_id);


--
-- Name: idx_sessions_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sessions_created_at ON public.sessions USING btree (created_at);


--
-- Name: idx_sessions_last_used; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sessions_last_used ON public.sessions USING btree (last_used);


--
-- Name: idx_task_runs_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_task_runs_created_at ON public.task_runs USING btree (created_at);


--
-- Name: idx_task_runs_job_run_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_task_runs_job_run_id ON public.task_runs USING btree (job_run_id);


--
-- Name: idx_task_runs_task_spec_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_task_runs_task_spec_id ON public.task_runs USING btree (task_spec_id);


--
-- Name: idx_task_specs_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_task_specs_deleted_at ON public.task_specs USING btree (deleted_at);


--
-- Name: idx_task_specs_job_spec_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_task_specs_job_spec_id ON public.task_specs USING btree (job_spec_id);


--
-- Name: idx_task_specs_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_task_specs_type ON public.task_specs USING btree (type);


--
-- Name: idx_tx_attempts_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_tx_attempts_created_at ON public.tx_attempts USING btree (created_at);


--
-- Name: idx_tx_attempts_tx_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_tx_attempts_tx_id ON public.tx_attempts USING btree (tx_id);


--
-- Name: idx_txes_from; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_txes_from ON public.txes USING btree ("from");


--
-- Name: idx_txes_nonce; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_txes_nonce ON public.txes USING btree (nonce);


--
-- Name: idx_users_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_created_at ON public.users USING btree (created_at);


--
-- Name: initiators initiators_job_spec_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.initiators
    ADD CONSTRAINT initiators_job_spec_id_fkey FOREIGN KEY (job_spec_id) REFERENCES public.job_specs(id);


--
-- Name: job_runs job_runs_job_spec_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.job_runs
    ADD CONSTRAINT job_runs_job_spec_id_fkey FOREIGN KEY (job_spec_id) REFERENCES public.job_specs(id);


--
-- Name: service_agreements service_agreements_job_spec_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.service_agreements
    ADD CONSTRAINT service_agreements_job_spec_id_fkey FOREIGN KEY (job_spec_id) REFERENCES public.job_specs(id);


--
-- Name: task_runs task_runs_job_run_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.task_runs
    ADD CONSTRAINT task_runs_job_run_id_fkey FOREIGN KEY (job_run_id) REFERENCES public.job_runs(id) ON DELETE CASCADE;


--
-- Name: task_specs task_specs_job_spec_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.task_specs
    ADD CONSTRAINT task_specs_job_spec_id_fkey FOREIGN KEY (job_spec_id) REFERENCES public.job_specs(id);


--
-- PostgreSQL database dump complete
--

