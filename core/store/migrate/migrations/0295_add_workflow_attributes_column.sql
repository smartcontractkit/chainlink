-- +goose Up
ALTER TABLE workflow_specs_v2 ADD COLUMN attributes bytea DEFAULT '';

-- +goose Down
ALTER TABLE workflow_specs_v2 DROP COLUMN attributes;
