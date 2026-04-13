# Postgres / txdb contention audit (manual)

## Scope

Repo-wide grep on `core/**/*_test.go` for patterns that commonly cause cross-test interference when many tests share one Postgres instance behind **txdb** (one transaction per test, rolled back at end).

## Findings (as of last audit)

| Pattern | Matches in `core/**/*_test.go` | Notes |
|--------|-------------------------------:|-------|
| `pg_advisory_lock` / `advisory_lock` | 0 | — |
| `CREATE TABLE` / `ALTER TABLE` / `CREATE INDEX` | 1 file | [txdb_test.go](./txdb_test.go) intentionally creates `txdb_test` inside a txdb transaction for driver behavior checks. |

Most DB tests use `NewSqlxDB` + DML only; DDL outside this helper is rare.

## When to use `heavyweight.FullTestDBV2` ([orm.go](../../utils/testutils/heavyweight/orm.go))

Prefer heavyweight when tests need:

- Multiple real connections seeing committed state
- Non-transactional Postgres behavior (certain DDL, extensions)
- Drivers other than txdb against a dedicated database

## Diagnostics

- `CHAINLINK_PGTEST_LOG_DB_CONCURRENCY=true` — logs when concurrent txdb sessions exceed 50 ([pgtest.go](./pgtest.go)).
- [PeakConcurrentTxDBSessions](./pgtest.go) — high-water concurrent `NewSqlxDB` count in-process.
- [LogPostgresActivitySummary](./activity_stats.go) — call from a test to print `pg_stat_activity` aggregates.
