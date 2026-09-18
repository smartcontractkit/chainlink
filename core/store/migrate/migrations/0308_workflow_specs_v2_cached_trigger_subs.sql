-- +goose Up

-- Trigger subscription payloads returned by a workflow's WASM Subscribe() call.
-- workflow_id is content-addressed (hash of owner+name+binary+config), so a
-- given row's payload never needs to change in place; cached once, reused on
-- every future engine start until the row itself is replaced/deleted.
ALTER TABLE workflow_specs_v2 ADD COLUMN trigger_subscriptions bytea;

-- +goose Down

ALTER TABLE workflow_specs_v2 DROP COLUMN trigger_subscriptions;
