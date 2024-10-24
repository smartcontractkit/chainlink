-- +goose Up
-- +goose StatementBegin

CREATE TABLE llo_retirement_report_cache (
    config_digest BYTEA NOT NULL CHECK (OCTET_LENGTH(config_digest) = 32) PRIMARY KEY,
    attested_retirement_report BYTEA NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE llo_retirement_report_cache_configs (
    config_digest BYTEA CHECK (octet_length(config_digest) = 32) PRIMARY KEY,
    signers BYTEA[] NOT NULL,
    f INT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE llo_retirement_report_cache_configs;
DROP TABLE llo_retirement_report_cache;

-- +goose StatementEnd
