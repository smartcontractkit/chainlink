-- +goose Up
ALTER TABLE upkeep_registrations
    ADD COLUMN IF NOT EXISTS last_keeper_index integer DEFAULT NULL;
ALTER TABLE upkeep_registrations
    DROP COLUMN IF EXISTS positioning_constant;
ALTER TABLE keeper_registries
    ADD COLUMN IF NOT EXISTS keeper_index_map jsonb DEFAULT NULL;

-- +goose Down
ALTER TABLE upkeep_registrations
    DROP COLUMN IF EXISTS last_keeper_index;
ALTER TABLE upkeep_registrations
    ADD COLUMN IF NOT EXISTS positioning_constant integer;
ALTER TABLE keeper_registries
    DROP COLUMN IF EXISTS keeper_index_map;
