-- +goose Up
CREATE TABLE dkg_shares(
    config_digest bytea NOT NULL CHECK ( length(config_digest) = 32 ),
    key_id bytea NOT NULL CHECK ( length(key_id) = 32 ),
    dealer bytea NOT NULL CHECK ( length(dealer) = 1),
    marshaled_share_record bytea NOT NULL,
    record_hash bytea NOT NULL CHECK ( length(record_hash) = 32 ),
    PRIMARY KEY (config_digest, key_id, dealer)
);

-- +goose Down
DROP TABLE dkg_shares;
