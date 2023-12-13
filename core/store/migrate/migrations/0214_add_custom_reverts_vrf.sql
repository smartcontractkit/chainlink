-- +goose Up
ALTER TABLE vrf_specs ADD COLUMN custom_reverts_pipeline_enabled boolean DEFAULT FALSE NOT NULL;

-- +goose Down
ALTER TABLE vrf_specs DROP COLUMN custom_reverts_pipeline_enabled;
