-- +goose Up
ALTER TABLE vrf_specs
    RENAME COLUMN confirmations TO min_incoming_confirmations;

-- +goose Down
ALTER TABLE vrf_specs
    RENAME COLUMN min_incoming_confirmations TO confirmations;