---
name: fix-flaky-go-test
description: >-
  Fix flaky Go tests in Chainlink: stress, Postgres, -shuffle, race (tools/bin),
  build tags. Use for intermittent failures, CI-only, -count/-shuffle issues,
  races, noisy output.
---

# Fix flaky Go tests (Chainlink)

<scope>
Reproduce before refactors. Fix determinism, isolation, time, concurrency.
Do not widen assertions or add blind retries.
Core tests need Postgres and usually CL_DATABASE_URL. CI uses tools/bin (gotestsum, race, integration), not only go test ./...
Read README.md Running tests, .github/workflows/ci-core.yml, tools/bin for parity.
</scope>

<setup>
Run README prep: pnpm, make mockery, make generate, Postgres, make setup-testdb, source .dbenv, make testdb after pulls. Use make testdb-force if DB stuck.
Unset env vars except CL_DATABASE_URL when tests act wrong.
CL_DATABASE_URL must target a *_test database (preparetest).
Modules: repo root, integration-tests/, core/scripts/. Run go test from the correct module root.
</setup>

<requirements>
If unknown, ask: package path, test name, module root, whether file is //go:build integration, whether test uses pgtest/cltest/SqlxDB or is -short safe.
State your assumptions when you start.
</requirements>

<principles>
Stress with plain go test -count/-failfast/-shuffle; gotestsum --rerun-fails in tools/bin/go_core_tests can hide flakes on PRs.
Treat flakes as production bugs until disproved.
Prefer injected time, IO, randomness; per-test resources; scoped state.
Do not loosen timeouts or assertions without a named cause.
</principles>

<classify>
Append --tags integration to every go test below if the file has //go:build integration.
deployment/ CCIP: use tools/bin/go_core_ccip_deployment_tests pattern (cd deployment, CL_RESERVE_PORTS=128).
Optional CI parity: GODEBUG=goindex=0 on go test (see ci-core.yml).
If the file uses //go:build dev or trace, add matching --tags when reproducing.
</classify>

<workflow>
<reproduce>
Stop when you have a stable repro. Add -v when needed.
Record package, -run regex, failure mode.

1. No DB quick path:
```sh
go test -short ./path/to/pkg -run '^TestName$' -count 100 -failfast
```

2. With DB from repo root:
```sh
source .dbenv && make testdb
go test ./path/to/pkg -run '^TestName$' -count 100 -failfast
```

3. Whole package: same DB prep then go test ./path/to/pkg -count 100 -failfast

4. Shuffle: add -shuffle on; bisect with -shuffle N

5. Race (fail if race.* exists):
```sh
GORACE="log_path=$PWD/race" go test -race -shuffle on -timeout 10s -count 100 ./path/to/pkg -run '^TestName$' -failfast
```

6. Parallelism probe: -cpu 1,2,4 and -parallel 4 with -shuffle on -count 50 -failfast

7. Optional full unit job after local repro: GODEBUG=goindex=0 ./tools/bin/go_core_tests ./... (see script for GITHUB_EVENT_NAME flags)
</reproduce>

<fix>
Apply fix_patterns. Avoid permanent time.Sleep as the main fix.
Re-run the same repro command. Record shuffle seed in commit or comment if order-dependent.
</fix>
</workflow>

<root_causes>
General: package init and globals, t.Parallel plus shared fixtures, wall clock without fakes, port or path collisions, map order assumptions, leaked env or cwd, goroutines after test end.

Chainlink: shared Postgres or stale schema; missing pgtest.NewSqlxDB(t); cltest.TestApplication teardown or leaked HTTP; ports without :0 or CL_RESERVE_PORTS; stress without --tags integration on integration files; wrong module root.
</root_causes>

<fix_patterns>
Scope state per test. Use t.Cleanup only when needed and obvious. Inject time, randomness, net, fs. Use t.TempDir and :0 listeners. Serialize or drop t.Parallel on shared resources. Prefer channels, WaitGroup, explicit sync over sleep polls.

Chainlink: pgtest.NewSqlxDB(t) and core/internal/testutils/pgtest helpers; testutils.Context(t); core/internal/cltest TestApplication and matching cleanup; configtest and evmtest under core/internal/testutils; core/utils/testutils/heavyweight for ORM-heavy tests.
</fix_patterns>

<verify>
Write the exact repro go test line including -run and --tags integration when relevant.
Race: GORACE log_path, go test -race -shuffle on, confirm no race.* or document skip.
Optional: TIMEOUT and COUNT with ./tools/bin/go_core_race_tests.
Do not merge unexplained timeout or assertion loosening.
</verify>
