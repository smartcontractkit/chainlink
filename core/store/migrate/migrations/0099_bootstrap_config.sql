-- +goose Up
-- +goose StatementBegin
CREATE TABLE bootstrap_contract_configs
(
    bootstrap_spec_id       INTEGER PRIMARY KEY,
    config_digest           bytea                    NOT NULL,
    config_count            bigint                   NOT NULL,
    signers                 bytea[],
    transmitters            text[],
    f                       smallint                 NOT NULL,
    onchain_config          bytea,
    offchain_config_version bigint                   NOT NULL,
    offchain_config         bytea,
    created_at              timestamp with time zone NOT NULL,
    updated_at              timestamp with time zone NOT NULL,
    CONSTRAINT bootstrap_contract_configs_config_digest_check CHECK ((octet_length(config_digest) = 32))
);

ALTER TABLE ONLY bootstrap_contract_configs
    ADD CONSTRAINT bootstrap_contract_configs_oracle_spec_fkey
        FOREIGN KEY (bootstrap_spec_id)
            REFERENCES bootstrap_specs (id)
            ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bootstrap_contract_configs;
-- +goose StatementEnd