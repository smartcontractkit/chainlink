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

<trunk>
## Trunk.io — gather context before touching code

Trunk.io tracks flaky test history, failure rates, and AI-generated root cause analysis.
Always check Trunk first — it may already have a fix recommendation.

### Finding the Trunk test link

Jira tickets for flaky tests almost always contain a Trunk link. Look for:
- A URL matching `https://app.trunk.io/chainlink/flaky-tests/test/<uuid>/`
- The UUID in that URL is the **test case ID** (not a fix/investigation ID)

To extract it from a Jira ticket:
```
mcp__atlassian__getJiraIssue  issue: "CCIP-XXXX"
```
Then look for `app.trunk.io` URLs in the description or comments.

### Reading test history and failure data

Open the test case page directly — it shows failure rate, timeline, and recent CI runs:
```
https://app.trunk.io/chainlink/flaky-tests/test/<test-case-uuid>/
```

Use the Scrapling MCP to fetch it (JS-rendered page):
```
mcp__ScraplingServer__fetch  url: "https://app.trunk.io/chainlink/flaky-tests/test/<uuid>/"
```

### Getting an AI fix recommendation (Trunk MCP)

The Trunk MCP tool `fix-flaky-test` requires a **fix/investigation ID**, which is different
from the test case ID in the URL. Investigations must be triggered from the Trunk UI first.

If an investigation exists, call:
```
mcp__plugin_trunk_trunk__fix-flaky-test
  repoName: "smartcontractkit/chainlink"
  orgSlug:  "chainlink"
  fixId:    "<fix-or-investigation-uuid>"
```

If the tool returns "Investigation not found", the investigation has not been triggered yet.
Ask the reporter to open the test case page and click "Investigate" — or proceed with
code-level analysis using the workflow below.

### What to read from Trunk

- **Failure rate** — how often it fails (e.g. 12% over last 30 days)
- **Failure pattern** — does it cluster around certain times, branches, or PR authors?
- **First seen / last seen** — did it regress recently after a change?
- **CI job name** — which workflow step fails (unit, race, integration, ccip-deployment)
- **Trunk root cause label** — if already classified (race, timing, docker, network, etc.)
</trunk>

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

Docker/Solana: WithSolanaContainerN port conflicts or slow startup; sync.Once download helpers that mark a failed download as done (causing cascading file-not-found failures in parallel runs); LoadCCIPPrograms network timeouts. If a test spins up Docker/Solana but the code under test early-exits before any chain interaction, remove the unnecessary infra.
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
