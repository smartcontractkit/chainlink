-- +goose Up
ALTER TABLE vrf_specs ADD COLUMN chunk_size bigint NOT NULL DEFAULT 20;

-- +goose Down
ALTER TABLE vrf_specs DROP COLUMN chunk_size;
