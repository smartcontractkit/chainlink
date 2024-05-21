-- +goose Up
ALTER TABLE vrf_specs
    ADD COLUMN "requested_confs_delay" BIGINT CHECK (requested_confs_delay >= 0) DEFAULT 0 NOT NULL;

-- +goose Down
ALTER TABLE vrf_specs
    DROP COLUMN "requested_confs_delay";
