-- +goose Up
CREATE SCHEMA evm;
SET search_path TO evm,public;

ALTER TABLE public.evm_forwarders SET SCHEMA evm;
ALTER TABLE evm.evm_forwarders RENAME TO forwarders;

ALTER TABLE public.evm_heads SET SCHEMA evm;
ALTER TABLE evm.evm_heads RENAME TO heads;

ALTER TABLE public.evm_key_states SET SCHEMA  evm;
ALTER TABLE evm.evm_key_states RENAME TO key_states;

ALTER TABLE public.evm_log_poller_blocks SET SCHEMA  evm;
ALTER TABLE evm.evm_log_poller_blocks RENAME TO log_poller_blocks;

ALTER TABLE public.evm_log_poller_filters SET SCHEMA evm;
ALTER TABLE evm.evm_log_poller_filters RENAME TO log_poller_filters;



ALTER TABLE public.evm_upkeep_states SET SCHEMA  evm;
ALTER TABLE evm.evm_upkeep_states RENAME TO upkeep_states;

ALTER TABLE public.eth_receipts SET SCHEMA  evm;
ALTER TABLE evm.eth_receipts RENAME TO  receipts;

ALTER TABLE public.eth_tx_attempts  SET SCHEMA evm;
ALTER TABLE evm.eth_tx_attempts  RENAME TO tx_attempts;

---------------------
-- Handle log triggers 
---------------------
DROP TRIGGER IF EXISTS notify_insert_on_evm_logs_topics ON PUBLIC.evm_logs;
DROP FUNCTION IF EXISTS public.notifysavedlogtopics();

ALTER TABLE public.evm_logs SET SCHEMA evm;
ALTER TABLE evm.evm_logs RENAME TO logs;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION evm.notifysavedlogtopics() RETURNS trigger
    LANGUAGE plpgsql
AS $$
BEGIN
    PERFORM pg_notify(
        'evm.insert_on_logs'::text,
        -- hex encoded address plus comma separated list of hex encoded topic values
        -- e.g. "<address>:<topicVal1>,<topicVal2>"
        encode(NEW.address, 'hex') || ':' || array_to_string(array(SELECT encode(unnest(NEW.topics), 'hex')), ',')
    );
    RETURN NULL;
END
$$;

DROP TRIGGER IF EXISTS notify_insert_on_logs_topics ON evm.logs;
CREATE TRIGGER notify_insert_on_logs_topics AFTER INSERT ON evm.logs FOR EACH ROW EXECUTE PROCEDURE evm.notifysavedlogtopics();
-- +goose StatementEnd

---------------------
-- Handle tx triggers
---------------------
DROP TRIGGER IF EXISTS notify_eth_tx_insertion on public.eth_txes; 
DROP FUNCTION IF EXISTS public.notifyethtxinsertion();

ALTER TABLE public.eth_txes SET SCHEMA  evm;
ALTER TABLE evm.eth_txes RENAME TO  txes;


-- +goose StatementBegin
CREATE OR REPLACE FUNCTION evm.notifytxinsertion() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
        BEGIN
		PERFORM pg_notify('evm.insert_on_txes'::text, encode(NEW.from_address, 'hex'));
		RETURN NULL;
        END
        $$;

DROP TRIGGER IF EXISTS notify_tx_insertion on evm.txes;
CREATE TRIGGER notify_tx_insertion AFTER INSERT ON evm.txes FOR EACH ROW EXECUTE PROCEDURE evm.notifytxinsertion();
-- +goose StatementEnd


-- +goose Down
SET search_path TO evm,public;
ALTER TABLE evm.forwarders SET SCHEMA public;
ALTER TABLE public.forwarders RENAME TO evm_forwarders;

ALTER TABLE evm.heads SET SCHEMA public;
ALTER TABLE public.heads RENAME TO evm_heads;

ALTER TABLE evm.key_states SET SCHEMA  public;
ALTER TABLE public.key_states RENAME TO evm_key_states;

ALTER TABLE evm.log_poller_blocks SET SCHEMA  public;
ALTER TABLE public.log_poller_blocks RENAME TO evm_log_poller_blocks;

ALTER TABLE evm.log_poller_filters SET SCHEMA public;
ALTER TABLE public.log_poller_filters RENAME TO evm_log_poller_filters;

ALTER TABLE evm.upkeep_states SET SCHEMA  public;
ALTER table public.upkeep_states RENAME TO evm_upkeep_states;

ALTER TABLE evm.receipts SET SCHEMA  public;
ALTER TABLE public.receipts RENAME TO eth_receipts;

ALTER TABLE evm.tx_attempts  SET SCHEMA public;
ALTER TABLE public.tx_attempts  RENAME TO eth_tx_attempts;


---------------------
-- Handle log triggers
---------------------

DROP TRIGGER IF EXISTS notify_insert_on_logs_topics ON evm.logs;
DROP FUNCTION IF EXISTS evm.notifysavedlogtopics();

ALTER TABLE evm.logs SET SCHEMA public;
ALTER TABLE public.logs RENAME TO evm_logs;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION PUBLIC.notifysavedlogtopics() RETURNS trigger
    LANGUAGE plpgsql
AS $$
BEGIN
    PERFORM pg_notify(
        'insert_on_evm_logs'::text,
        -- hex encoded address plus comma separated list of hex encoded topic values
        -- e.g. "<address>:<topicVal1>,<topicVal2>"
        encode(NEW.address, 'hex') || ':' || array_to_string(array(SELECT encode(unnest(NEW.topics), 'hex')), ',')
    );
    RETURN NULL;
END
$$;

DROP TRIGGER IF EXISTS notify_insert_on_evm_logs_topics ON PUBLIC.evm_logs;
CREATE TRIGGER notify_insert_on_evm_logs_topics AFTER INSERT ON PUBLIC.evm_logs FOR EACH ROW EXECUTE PROCEDURE PUBLIC.notifysavedlogtopics();
-- +goose StatementEnd

---------------------
-- Handle tx triggers
---------------------

DROP TRIGGER IF EXISTS notify_tx_insertion on evm.txes; 
DROP FUNCTION IF EXISTS evm.notifytxinsertion();

ALTER TABLE evm.txes SET SCHEMA  public;
ALTER TABLE public.txes RENAME TO eth_txes;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION public.notifyethtxinsertion() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
        BEGIN
		PERFORM pg_notify('insert_on_eth_txes'::text, encode(NEW.from_address, 'hex'));
		RETURN NULL;
        END
        $$;

DROP TRIGGER IF EXISTS notify_eth_tx_insertion on public.eth_txes;
CREATE TRIGGER notify_eth_tx_insertion AFTER INSERT ON public.eth_txes FOR EACH ROW EXECUTE PROCEDURE public.notifyethtxinsertion();
-- +goose StatementEnd

DROP SCHEMA evm;
SET search_path TO public;