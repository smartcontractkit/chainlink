-- +goose Up
ALTER TABLE external_initiators ADD CONSTRAINT external_initiators_name_unique UNIQUE(name);
ALTER TABLE webhook_specs ADD COLUMN external_initiator_name TEXT REFERENCES external_initiators (name);
ALTER TABLE webhook_specs ADD COLUMN external_initiator_spec JSONB;
ALTER TABLE webhook_specs ADD CONSTRAINT external_initiator_null_not_null CHECK (
    external_initiator_name IS NULL AND external_initiator_spec IS NULL
        OR
    external_initiator_name IS NOT NULL AND external_initiator_spec IS NOT NULL
);

-- +goose Down
ALTER TABLE external_initiators DROP CONSTRAINT external_initiators_name_unique;
ALTER TABLE webhook_specs DROP COLUMN external_initiator_name;
ALTER TABLE webhook_specs DROP COLUMN external_initiator_spec;
ALTER TABLE webhook_specs DROP CONSTRAINT external_initiator_null_not_null;
