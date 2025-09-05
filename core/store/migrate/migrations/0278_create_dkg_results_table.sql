-- +goose Up
CREATE TABLE dkg_results (
    instance_id TEXT PRIMARY KEY,
    config_digest BYTEA NOT NULL,
    seq_nr BIGINT NOT NULL,
    report_with_result_package BYTEA NOT NULL,
    signatures BYTEA[] NOT NULL,
    signer_oracle_ids BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

-- +goose Down
DROP TABLE dkg_results;
