-- +goose Up
-- +goose StatementBegin

ALTER TABLE cre_settings_specs
    ADD COLUMN config_type TEXT NOT NULL DEFAULT 'settings',
    ADD COLUMN shard_assignment TEXT;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE cre_settings_specs
    DROP COLUMN shard_assignment,
    DROP COLUMN config_type;

-- +goose StatementEnd
