-- +goose Up
-- config_type selects how a cresettings payload is interpreted; offchain_config carries the
-- proto-JSON payload for config_type = 'capabilities_registry'. Both default to '' so existing
-- rows keep resolving exactly as before (config_type embedded in settings, else "settings").
ALTER TABLE cre_settings_specs
    ADD COLUMN config_type TEXT NOT NULL DEFAULT '',
    ADD COLUMN offchain_config TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE cre_settings_specs
    DROP COLUMN offchain_config,
    DROP COLUMN config_type;
