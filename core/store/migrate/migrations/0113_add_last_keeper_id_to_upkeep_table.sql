-- +goose Up
ALTER TABLE upkeep_registrations
    ADD COLUMN last_keeper_index integer DEFAULT NULL;
ALTER TABLE upkeep_registrations
    DROP COLUMN positioning_constant;


-- +goose Down
ALTER TABLE upkeep_registrations
    DROP COLUMN last_keeper_index;
ALTER TABLE upkeep_registrations
    ADD COLUMN positioning_constant integer;
