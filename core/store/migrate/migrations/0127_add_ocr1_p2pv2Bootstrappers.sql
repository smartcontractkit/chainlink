-- +goose Up
ALTER TABLE ocr_oracle_specs
    ADD COLUMN p2pv2_bootstrappers text[] NOT NULL DEFAULT '{}';

-- +goose Down
ALTER TABLE ocr_oracle_specs
    DROP COLUMN p2pv2_bootstrappers;
