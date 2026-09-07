-- +goose Up
ALTER TABLE workflow_specs_v2 ADD COLUMN registered_at bigint NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE workflow_specs_v2 DROP COLUMN registered_at;
