-- +goose Up
ALTER TABLE upkeep_registrations
    ADD COLUMN IF NOT EXISTS last_keeper_index integer DEFAULT NULL;
ALTER TABLE keeper_registries
    ADD COLUMN IF NOT EXISTS keeper_index_map jsonb DEFAULT NULL;

-- +goose Down
ALTER TABLE upkeep_registrations
    DROP COLUMN IF EXISTS last_keeper_index;
ALTER TABLE keeper_registries
    DROP COLUMN IF EXISTS keeper_index_map;
