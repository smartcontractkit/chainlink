-- +goose Up
UPDATE jobs SET forwarding_allowed = false WHERE forwarding_allowed is NULL;
ALTER TABLE jobs ALTER COLUMN forwarding_allowed SET NOT NULL;

-- +goose Down
ALTER TABLE jobs ALTER COLUMN forwarding_allowed DROP NOT NULL;
