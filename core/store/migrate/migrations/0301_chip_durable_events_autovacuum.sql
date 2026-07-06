-- +goose Up

-- More aggressive per-table autovacuum for the high-churn durable-events queue.
-- Helps prevent against slow sql errors.
-- Under delete-on-delivery every event is INSERTed then DELETEd, so the table
-- generates dead tuples at roughly the event rate (thousands/s under load).
ALTER TABLE cre.chip_durable_events SET (
    autovacuum_vacuum_scale_factor  = 0.0,
    autovacuum_vacuum_threshold     = 10000,
    autovacuum_vacuum_cost_delay    = 1,
    autovacuum_vacuum_cost_limit    = 10000,
    autovacuum_analyze_scale_factor = 0.0,
    autovacuum_analyze_threshold    = 10000
);

-- +goose Down

ALTER TABLE cre.chip_durable_events RESET (
    autovacuum_vacuum_scale_factor,
    autovacuum_vacuum_threshold,
    autovacuum_vacuum_cost_delay,
    autovacuum_vacuum_cost_limit,
    autovacuum_analyze_scale_factor,
    autovacuum_analyze_threshold
);
