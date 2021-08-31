-- +goose Up
CREATE TABLE offchainreporting_latest_round_requested (
	offchainreporting_oracle_spec_id integer PRIMARY KEY REFERENCES offchainreporting_oracle_specs (id) DEFERRABLE INITIALLY IMMEDIATE,
	requester bytea not null CHECK (octet_length(requester) = 20),
	config_digest bytea not null CHECK (octet_length(config_digest) = 16),
	epoch bigint not null,
	round bigint not null,
	raw jsonb not null
);

-- +goose Down
DROP TABLE offchainreporting_latest_round_requested;
