-- +goose Up
CREATE SCHEMA evm;
SET search_path TO evm,public;

ALTER TABLE public.evm_forwarders set schema evm;
ALTER TABLE evm.evm_forwarders RENAME TO forwarders;

ALTER TABLE public.evm_heads set schema evm;
ALTER TABLE evm.evm_heads RENAME TO heads;

ALTER TABLE public.evm_key_states set schema  evm;
ALTER TABLE evm.evm_key_states RENAME TO key_states;

ALTER TABLE public.evm_log_poller_blocks set schema  evm;
ALTER TABLE evm.evm_log_poller_blocks RENAME TO log_poller_blocks;

ALTER TABLE public.evm_log_poller_filters set schema evm;
ALTER TABLE evm.evm_log_poller_filters RENAME TO log_poller_filters;

ALTER TABLE public.evm_logs set schema evm;
ALTER TABLE evm.evm_logs RENAME TO logs;

ALTER TABLE public.evm_upkeep_states set schema  evm;
ALTER TABLE evm.evm_upkeep_states RENAME TO upkeep_states;

ALTER TABLE eth_receipts set schema  evm;
ALTER TABLE eth_tx_attempts  set schema evm;
ALTER TABLE eth_txes set schema  evm;

-- +goose Down
SET search_path TO evm,public;
ALTER TABLE evm.forwarders set schema public;
ALTER TABLE public.forwarders rename to evm_forwarders;

ALTER TABLE evm.heads set schema public;
ALTER TABLE public.heads RENAME TO evm_heads;

ALTER TABLE evm.key_states set schema  public;
ALTER TABLE public.key_states RENAME TO evm_key_states;

ALTER TABLE evm.log_poller_blocks set schema  public;
ALTER TABLE public.log_poller_blocks RENAME TO evm_log_poller_blocks;

ALTER TABLE evm.log_poller_filters set schema public;
ALTER TABLE public.log_poller_filters RENAME TO evm_log_poller_filters;

ALTER TABLE evm.logs set schema public;
ALTER TABLE public.logs RENAME TO evm_logs;

ALTER TABLE evm.upkeep_states set schema  public;
ALTER table public.upkeep_states RENAME TO evm_upkeep_states;

ALTER TABLE eth_receipts set schema  public;
ALTER TABLE eth_tx_attempts  set schema public;
ALTER TABLE eth_txes set schema  public;

DROP SCHEMA evm;
