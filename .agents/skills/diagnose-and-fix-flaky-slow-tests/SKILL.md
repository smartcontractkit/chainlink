---
name: diagnose-and-fix-flaky-slow-tests
description: >-
  Runs the Chainlink `tools/test` survey loop to reproduce flaky or slow Go
  tests, reads `report.json` and per-iteration JSONL logs, then fixes root causes.
  Use when investigating flaky tests, intermittent failures, test timeouts,
  races, or tests that exceed a time budget; when the user mentions survey,
  reruns, flakes, or slow tests in this repo.
---

# Diagnose and fix flaky / slow tests (survey)

## Preconditions

- Work from the **Chainlink repo root** (the harness discovers `repo_root` from the working directory).
- Invoke the harness exactly as documented in `tools/test/AGENTS.md`:

```sh
go -C ./tools/test run . survey [survey flags] <single go test target>
```

- **One target only** (e.g. `./core/...`, `./pkg/foo/...`). Additional packages or multiple args are rejected.
- Each iteration runs **`go test -json -count=1`** and writes JSONL.

## What survey does

1. Creates `test-survey-results-<YYYYMMDDhhmmss>/` under the repo root.
2. For each iteration, runs `go test -json` for the target with `-timeout` from `--timeout`, capturing stdout/stderr into `iteration-<n>.log.jsonl`.
3. After all iterations (or interrupt), parses those streams and writes **`report.json`** (pretty JSON) in the same results directory.
4. Between iterations **after the first**, the test database is reset when using the default ephemeral Postgres setup—useful for stateful flakes.

Interrupt (`SIGINT`/`SIGTERM`) stops the loop but still triggers analysis on completed `iteration-*.log.jsonl` files.

## Flags that matter

| Flag                 | Role                                                                                                                                                |
| -------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------- |
| `--iterations N`     | Repeat full target run N times (default 1). Use **≥10** for flake hunting when runtime allows; **3–5** for a quick smoke.                           |
| `--timeout`          | Per-iteration `go test -timeout` (default 10m). Raise for slow suites; lower to surface hangs faster.                                               |
| `--slow-threshold`   | Tests whose **max** elapsed across iterations exceeds this are listed under `slow` in `report.json` (default 30s). Set `0` to disable slow listing. |
| `--fail-fast`        | Stop after the first failing iteration (good for tight edit/rerun loops; **bad** for flake statistics).                                             |
| `--ai-output`        | Sparse, agent-friendly stdout; use this when running                                                                                                |
| `--database-url`     | Optional: use an existing Postgres instead of ephemeral.                                                                                            |
| `--postgres-version` | Postgres version for the ephemeral DB (default in `tools/test/internal/config`).                                                                    |

## Agent workflow

1. **Shrink the target** when possible: survey accepts **exactly one** positional argument (the `go test` package pattern, e.g. `./core/foo/...`). Prefer the smallest `./path/...` that contains the suspected test; there is no separate hook to pass `-run` without changing `tools/test`. Locate a specific test name inside `report.json`, `iteration-*.log.jsonl`, or by searching log output.

2. **Run survey** with iterations suited to suspicion:
   - Quick check: `--iterations 5`
   - Flake hunt: `--iterations 10`–`30`, aligned with CI budget

   Example:

   ```sh
   go -C ./tools/test run . survey --iterations 15 --timeout=15m ./core/services/somepkg/...
   ```

3. **Open the newest results directory** (by timestamp suffix) and read **`report.json`**.

4. **Interpret `report.json` sections** (see `tools/test/internal/runner/analyze.go`):
   - **`flakes`**: same test had both pass and fail events across iterations → classic flake.
   - **`failures`**: failed every time (`passes == 0`) → deterministic failure (or environment).
   - **`timeouts`**: `panic: test timed out` attributed (with reattribution from `running tests:` when applicable).
   - **`slow`**: max elapsed across iterations > `slow_threshold` (and not a timeout bucket).

   Each relevant entry may include **`logs`**: per-failing-iteration captured output—use **`iteration`** indices to match **`iteration-<n>.log.jsonl`** for full JSONL streams when needed.

5. **Triage fixes** (typical patterns):
   - **Order-dependent / shared state**: isolate globals, reset fixtures, avoid parallel collisions; check `-parallel` and `t.Parallel()` interactions.
   - **Timing**: replace sleeps with synchronization (`require.Eventually`, channels, mocks); fix assumptions about wall clock.
   - **DB / migrations**: leverage the fact that DB resets between iterations—if flake disappears, suspect migration or connection pool hygiene.
   - **Timeouts**: reduce work, split tests, or raise only where justified; fix deadlocks before bumping `-timeout`.

6. **Verify**: re-run survey on the same target with the same or higher `--iterations` until `flakes` for the target is empty (or acceptable). Run package/unit tests as usual per `tools/test/AGENTS.md`.

## Output locations (single source of truth)

```
test-survey-results-<timestamp>/
  iteration-0.log.jsonl
  iteration-1.log.jsonl
  ...
  report.json
```

## Commands for harness changes

If editing `tools/test`:

```sh
golangci-lint run ./... --fix
go test ./...
```

(from `tools/test` module root as appropriate).
