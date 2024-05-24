-- +goose Up
ALTER TABLE offchainreporting_oracle_specs ADD COLUMN database_timeout BIGINT;
ALTER TABLE offchainreporting_oracle_specs ADD COLUMN observation_grace_period BIGINT;
ALTER TABLE offchainreporting_oracle_specs ADD COLUMN contract_transmitter_transmit_timeout BIGINT;

-- +goose Down
ALTER TABLE offchainreporting_oracle_specs DROP COLUMN database_timeout;
ALTER TABLE offchainreporting_oracle_specs DROP COLUMN observation_grace_period;
ALTER TABLE offchainreporting_oracle_specs DROP COLUMN contract_transmitter_transmit_timeout;
