-- +goose Up

-- SDK name discovered when the workflow binary is first compiled. Stored with
-- trigger_subscriptions so engines can start without compiling the binary.
ALTER TABLE workflow_specs_v2 ADD COLUMN sdk_name text NOT NULL DEFAULT '';

-- +goose Down

ALTER TABLE workflow_specs_v2 DROP COLUMN sdk_name;
