CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;
COMMENT ON EXTENSION citext IS 'data type for case-insensitive character strings';
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;
COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';
CREATE FUNCTION public.copy_task_run_confirmations() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
        BEGIN
          NEW.confirmations_new1562419039813 = NEW.confirmations;
          NEW."minimumConfirmations_new1562419039813" = NEW."minimumConfirmations";
          RETURN NEW;
        END;
      $$;
SET default_tablespace = '';
SET default_with_oids = false;
CREATE TABLE public.chainlink_node (
    id bigint NOT NULL,
    "createdAt" timestamp without time zone DEFAULT now() NOT NULL,
    name character varying NOT NULL,
    "accessKey" character varying(32) NOT NULL,
    "hashedSecret" character varying(64) NOT NULL,
    salt character varying(64) NOT NULL,
    url character varying
);
INSERT INTO public.chainlink_node VALUES (1, '2019-07-09 15:48:50.250417', 'NodeyMcNodeFace', 'u4HULe0pj5xPyuvv', '302df2b42ab313cb9b00fe0cca9932dacaaf09e662f2dca1be9c2ad2d927d5df', 'wZ02sJ8iZ6WffxXduxwzkCfOc3PS8BZJ', NULL);
CREATE SEQUENCE public.chainlink_node_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.chainlink_node_id_seq OWNED BY public.chainlink_node.id;
CREATE TABLE public.job_run (
    id character varying DEFAULT public.uuid_generate_v4() NOT NULL,
    "runId" public.citext NOT NULL,
    "jobId" public.citext NOT NULL,
    status character varying NOT NULL,
    error character varying,
    "createdAt" timestamp without time zone DEFAULT now() NOT NULL,
    "finishedAt" timestamp without time zone,
    type character varying NOT NULL,
    "requestId" public.citext,
    "txHash" public.citext,
    requester public.citext,
    "chainlinkNodeId" bigint NOT NULL
);
CREATE TABLE public.migrations (
    id integer NOT NULL,
    "timestamp" bigint NOT NULL,
    name character varying NOT NULL
);
CREATE SEQUENCE public.migrations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.migrations_id_seq OWNED BY public.migrations.id;
INSERT INTO public.migrations VALUES (1, 1557261237896, 'InitialMigration1557261237896');
INSERT INTO public.migrations VALUES (2, 1559910921273, 'AddConfirmationsToTaskRun1559910921273');
INSERT INTO public.migrations VALUES (3, 1562419039813, 'ConvertTaskRunConfirmationsToBigInt1562419039813');
INSERT INTO public.migrations VALUES (4, 1564009523000, 'AddUrlsToNodes1564009523000');
CREATE TABLE public.task_run (
    id bigint NOT NULL,
    "jobRunId" character varying NOT NULL,
    index integer NOT NULL,
    type character varying NOT NULL,
    status character varying NOT NULL,
    error character varying,
    "transactionHash" character varying,
    "transactionStatus" character varying,
    confirmations integer,
    "minimumConfirmations" integer,
    confirmations_new1562419039813 bigint,
    "minimumConfirmations_new1562419039813" bigint
);
CREATE SEQUENCE public.task_run_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.task_run_id_seq OWNED BY public.task_run.id;
ALTER TABLE ONLY public.chainlink_node ALTER COLUMN id SET DEFAULT nextval('public.chainlink_node_id_seq'::regclass);
ALTER TABLE ONLY public.migrations ALTER COLUMN id SET DEFAULT nextval('public.migrations_id_seq'::regclass);
ALTER TABLE ONLY public.task_run ALTER COLUMN id SET DEFAULT nextval('public.task_run_id_seq'::regclass);
ALTER TABLE ONLY public.migrations
    ADD CONSTRAINT "PK_8c82d7f526340ab734260ea46be" PRIMARY KEY (id);
ALTER TABLE ONLY public.chainlink_node
    ADD CONSTRAINT "chainlink_node_accessKey_key" UNIQUE ("accessKey");
ALTER TABLE ONLY public.chainlink_node
    ADD CONSTRAINT chainlink_node_name_key UNIQUE (name);
ALTER TABLE ONLY public.chainlink_node
    ADD CONSTRAINT chainlink_node_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.job_run
    ADD CONSTRAINT job_run_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.task_run
    ADD CONSTRAINT task_run_pkey PRIMARY KEY (id);
CREATE UNIQUE INDEX chainlink_node_access_key_idx ON public.chainlink_node USING btree ("accessKey");
CREATE UNIQUE INDEX job_run_chainlink_node_id_run_id_idx ON public.job_run USING btree ("chainlinkNodeId", "runId");
CREATE INDEX job_run_job_id_idx ON public.job_run USING btree ("jobId");
CREATE INDEX job_run_request_id_idx ON public.job_run USING btree ("requestId");
CREATE INDEX job_run_requester_idx ON public.job_run USING btree (requester);
CREATE INDEX job_run_tx_hash_idx ON public.job_run USING btree ("txHash");
CREATE INDEX task_run_index_idx ON public.task_run USING btree (index);
CREATE UNIQUE INDEX task_run_index_job_run_id_idx ON public.task_run USING btree (index, "jobRunId");
CREATE INDEX task_run_job_run_id_idx ON public.task_run USING btree ("jobRunId");
CREATE TRIGGER check_task_run_confirmations BEFORE INSERT OR UPDATE ON public.task_run FOR EACH ROW WHEN (((new.confirmations IS NOT NULL) OR (new."minimumConfirmations" IS NOT NULL))) EXECUTE PROCEDURE public.copy_task_run_confirmations();
ALTER TABLE ONLY public.job_run
    ADD CONSTRAINT "job_run_chainlinkNodeId_fkey" FOREIGN KEY ("chainlinkNodeId") REFERENCES public.chainlink_node(id);
ALTER TABLE ONLY public.task_run
    ADD CONSTRAINT "task_run_jobRunId_fkey" FOREIGN KEY ("jobRunId") REFERENCES public.job_run(id);
