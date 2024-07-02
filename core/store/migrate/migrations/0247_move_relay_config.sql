-- +goose Up
-- +goose StatementBegin
ALTER TABLE jobs
    ADD COLUMN relay_config JSONB NOT NULL DEFAULT '{}',
    ADD COLUMN relay text NOT NULL DEFAULT '';

UPDATE jobs
    SET relay_config = ocr2_oracle_specs.relay_config,
        relay = ocr2_oracle_specs.relay
            FROM ocr2_oracle_specs
                WHERE jobs.ocr2_oracle_spec_id = ocr2_oracle_specs.id;

ALTER TABLE ocr2_oracle_specs
    DROP COLUMN relay_config,
    DROP COLUMN relay;
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
ALTER TABLE ocr2_oracle_specs
    ADD COLUMN relay_config JSONB NOT NULL DEFAULT '{}',
    ADD COLUMN relay text NOT NULL;

UPDATE ocr2_oracle_specs
    SET relay_config = jobs.relay_config,
        relay = jobs.relay
        FROM jobs
            WHERE jobs.ocr2_oracle_spec_id = ocr2_oracle_specs.id;

ALTER TABLE jobs
    DROP COLUMN relay_config,
    DROP COLUMN relay;

-- +goose StatementEnd
