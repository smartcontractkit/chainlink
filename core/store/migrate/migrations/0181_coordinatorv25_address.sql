-- +goose Up

ALTER TABLE vrf_specs ADD COLUMN vrf_version TEXT;

-- +goose Down

ALTER TABLE vrf_specs DROP COLUMN vrf_version;
