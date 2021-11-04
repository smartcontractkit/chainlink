-- +goose Up
ALTER TABLE keeper_specs
    ADD COLUMN confirmations integer;

-- +goose Down
ALTER TABLE keeper_specs
    DROP COLUMN confirmations;
