-- +goose Up

ALTER TABLE cre.trigger_pending_events
    ADD COLUMN IF NOT EXISTS org_id TEXT NOT NULL DEFAULT '';

-- +goose Down

ALTER TABLE cre.trigger_pending_events
    DROP COLUMN IF EXISTS org_id;
