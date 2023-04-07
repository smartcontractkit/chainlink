-- +goose Up
DROP INDEX idx_eth_tx_attempts_hash;
CREATE UNIQUE INDEX idx_eth_tx_attempts_hash ON public.eth_tx_attempts USING btree (hash) WHERE state != 'fatal_error';

-- +goose Down
DROP INDEX idx_eth_tx_attempts_hash;
CREATE UNIQUE INDEX idx_eth_tx_attempts_hash ON public.eth_tx_attempts USING btree (hash);
