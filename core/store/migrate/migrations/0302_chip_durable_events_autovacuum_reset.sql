-- +goose Up

-- Revert the aggressive per-table autovacuum overrides added in 0301 so the
-- durable-events queue falls back to the cluster/global autovacuum defaults.
-- Nodes are already migrated past 0301, so removing that migration would neither
-- undo the applied ALTER nor keep goose history consistent; this forward
-- migration cleanly un-applies the overrides instead. RESET drops the per-table
-- settings, reverting to whatever the cluster globals are.
ALTER TABLE cre.chip_durable_events RESET (
    autovacuum_vacuum_scale_factor,
    autovacuum_vacuum_threshold,
    autovacuum_vacuum_cost_delay,
    autovacuum_vacuum_cost_limit,
    autovacuum_analyze_scale_factor,
    autovacuum_analyze_threshold
);

-- +goose Down

-- Re-apply the aggressive per-table autovacuum settings from 0301.
ALTER TABLE cre.chip_durable_events SET (
    autovacuum_vacuum_scale_factor  = 0.0,
    autovacuum_vacuum_threshold     = 10000,
    autovacuum_vacuum_cost_delay    = 1,
    autovacuum_vacuum_cost_limit    = 10000,
    autovacuum_analyze_scale_factor = 0.0,
    autovacuum_analyze_threshold    = 10000
);
