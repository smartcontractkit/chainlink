-- +goose Up
ALTER TABLE upkeep_registrations
    ADD COLUMN last_keeper_index integer DEFAULT NULL;
ALTER TABLE upkeep_registrations
    DROP COLUMN positioning_constant;
ALTER TABLE keeper_registries
    ADD COLUMN keeper_index_map jsonb DEFAULT NULL;

-- +goose Down
ALTER TABLE upkeep_registrations
    DROP COLUMN last_keeper_index;
ALTER TABLE upkeep_registrations
    ADD COLUMN positioning_constant integer;
ALTER TABLE keeper_registries
    DROP COLUMN IF EXISTS keeper_index_map;
