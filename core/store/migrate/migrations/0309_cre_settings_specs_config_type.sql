-- +goose Up
-- config_type selects how a cresettings payload is interpreted; offchain_config carries the
-- proto-JSON payload for config_type = 'capabilities_registry'. Both default to '' so existing
-- rows keep resolving exactly as before (config_type embedded in settings, else "settings").
ALTER TABLE cre_settings_specs
    ADD COLUMN config_type TEXT NOT NULL DEFAULT '',
    ADD COLUMN offchain_config TEXT NOT NULL DEFAULT '';

-- The top-level config_type column is the only accepted discriminator for
-- capabilities_registry (the legacy config_type key embedded in settings is rejected for it),
-- and only that form may carry an offchain payload. This keeps the unique index below covering
-- every row the node would route to the offchain capabilities registry.
ALTER TABLE cre_settings_specs
    ADD CONSTRAINT chk_cre_settings_specs_config_type
        CHECK (config_type IN ('', 'settings', 'shard_assignment', 'capabilities_registry')),
    ADD CONSTRAINT chk_cre_settings_specs_offchain_config
        CHECK ((config_type = 'capabilities_registry') = (offchain_config <> '')),
    ADD CONSTRAINT chk_cre_settings_specs_cap_registry_settings
        CHECK (config_type <> 'capabilities_registry' OR settings = '');

-- At most one capabilities_registry job may exist. Enforcing this at insert time (rather than
-- only in the delegate) keeps a rejected duplicate from being persisted, where on the next boot
-- it could start before (and displace) the job that was actually applied.
CREATE UNIQUE INDEX idx_cre_settings_specs_single_capabilities_registry
    ON cre_settings_specs (config_type)
    WHERE config_type = 'capabilities_registry';

-- Durable high-water mark of the offchain capabilities registry payload version. It is checked
-- and advanced in the same transaction that inserts a capabilities_registry spec, so a rolled
-- back insert leaves it unchanged, and it survives job deletion and node restarts: a payload
-- older than the newest one ever committed is rejected (re-creating that same payload is
-- allowed).
CREATE TABLE cre_offchain_registry_high_water (
    id      INT PRIMARY KEY CHECK (id = 1),
    version NUMERIC(20, 0) NOT NULL,
    hash    TEXT NOT NULL
);
INSERT INTO cre_offchain_registry_high_water (id, version, hash) VALUES (1, 0, '');

-- +goose Down
DROP TABLE IF EXISTS cre_offchain_registry_high_water;

DROP INDEX IF EXISTS idx_cre_settings_specs_single_capabilities_registry;

ALTER TABLE cre_settings_specs
    DROP CONSTRAINT IF EXISTS chk_cre_settings_specs_cap_registry_settings,
    DROP CONSTRAINT IF EXISTS chk_cre_settings_specs_offchain_config,
    DROP CONSTRAINT IF EXISTS chk_cre_settings_specs_config_type,
    DROP COLUMN offchain_config,
    DROP COLUMN config_type;
