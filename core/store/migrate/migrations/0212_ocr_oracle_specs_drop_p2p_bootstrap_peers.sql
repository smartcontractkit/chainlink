-- +goose Up
ALTER TABLE ocr_oracle_specs DROP COLUMN p2p_bootstrap_peers;

-- +goose Down
ALTER TABLE ocr_oracle_specs ADD COLUMN p2p_bootstrap_peers text[];
