-- +goose Up
ALTER TABLE ocr2_oracle_specs
    RENAME COLUMN p2p_bootstrap_peers to p2pv2_bootstrappers;

-- +goose Down
ALTER TABLE ocr2_oracle_specs
    RENAME COLUMN p2pv2_bootstrappers to p2p_bootstrap_peers;
