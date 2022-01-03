-- +goose Up
ALTER TABLE offchainreporting_oracle_specs DROP COLUMN p2p_peer_id;
ALTER TABLE offchainreporting2_oracle_specs DROP COLUMN p2p_peer_id;

-- +goose Down
ALTER TABLE offchainreporting_oracle_specs ADD COLUMN p2p_peer_id TEXT;
ALTER TABLE offchainreporting2_oracle_specs ADD COLUMN p2p_peer_id TEXT;
