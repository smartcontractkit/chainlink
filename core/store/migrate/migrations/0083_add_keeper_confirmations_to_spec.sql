-- +goose Up
ALTER TABLE keeper_specs
    ADD COLUMN min_incoming_confirmations integer;

-- +goose Down
ALTER TABLE keeper_specs
    DROP COLUMN min_incoming_confirmations;
