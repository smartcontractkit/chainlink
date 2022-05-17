-- +goose Up
ALTER TABLE ocr2_oracle_specs
    DROP COLUMN p2p_bootstrap_peers;

-- +goose Down
ALTER TABLE ocr2_oracle_specs
    ADD COLUMN p2p_bootstrap_peers text[] NOT NULL DEFAULT '{}';