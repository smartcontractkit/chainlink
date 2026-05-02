---
name: chainlink-test-diagnosis
description: >-
  Diagnoses and fixes unstable Chainlink Go tests (flakes, races, timeouts, deadlocks,
  slow runs) using the repository `tools/test` diagnose flow and report/log analysis.
  Use when debugging non-deterministic failures, CI-only instability, or test runtime;
  skip deterministic failures, routine test runs, and work unrelated to test stability.
---

<source_of_truth>
Authoritative CLI flags and behavior: `go -C tools/test run . diagnose -h` (from repository root).
</source_of_truth>

<purpose>
Root-cause analysis and fixes for unstable or slow Chainlink Go tests.
Flow:
1. Scope target and hypothesis.
2. Run `diagnose` on bounded target.
3. Analyze `report.json` and logs.
4. Apply playbook. Fix. Verify with same scope.
</purpose>

<activation>
Use for diagnosing and fixing flakes, timeouts, deadlocks, CI-only failures, and slow tests.
Skip for deterministic first-run failures, typical test runs, known fixes, full-suite CI prep, and unrelated test work.
</activation>

<package_scope>
Skill is designed to run tests in these packages.
Running all tests in high-level packages/directories can take a long time (6m+ each iteration).
Do not run before warning the user how long a full diagnose run could take and getting explicit approval.
- `core/`
- `deployment/`

If running tests not in those packages, warn the user the skill is not designed for it and might have unexpected outputs/breakage.
</package_scope>

<preflight>
- Ask for test, package, issue, or permission to discover.
- Start bounded: single test, package, or subtree; use `--fail-fast` or low `--iterations`.
- Classify hypothesis: flake, timeout, slow, panic, deadlock, or race.
- Put `go test` flags after `--`.
</preflight>

<diagnose_cli>
<restrictions>
Sandbox blocks `diagnose` when it needs local Postgres. Ask the user to run the command or approve unsandboxed execution if this error appears.
<error>
Failed to reset database:unable to drop postgres database: failed to connect to `host=localhost user=postgres database=template1`: dial error (dial tcp 127.0.0.1:55001: connect: operation not permitted)
</error>
</restrictions>

Run from the chainlink repository root. Use `go -C tools/test ...`. Equivalent Make targets: `make new_test_diagnose ARGS='...'`, `make new_test ARGS='...'`, `make new_gotestsum ARGS='...'`.

Harness-only flags before `--`: `--iterations`, `--parallel-iterations`, `--slow-threshold`, `--fail-fast`, `--fail-fast-on`, `--shuffle-seed`.
Everything after `--` passes to `go test`: `-timeout`, `-race`, `-run`, package patterns.
Put package patterns last.

```sh
# Command help
go -C tools/test run . diagnose -h
# Example: harness flags, then --, then go test flags and packages
go -C tools/test run . diagnose --iterations <N> --parallel-iterations <N> --slow-threshold <duration> --fail-fast --ai-output -- --timeout <duration> --run '<regex>' --race ./path/to/package/...
# Stop only when the hunted signal appears
go -C tools/test run . diagnose --iterations <N> --fail-fast-on=timeout -- --timeout <duration> ./path/to/package/...
go -C tools/test run . diagnose --iterations <N> --fail-fast-on=slow --slow-threshold <duration> -- ./path/to/package/...
```

<semantics>
- `diagnose` prepends `go test -json` each iteration and drops duplicate `-json`.
- `diagnose` adds `-count=1` unless `go test` sets `-count` greater than 1. Use diagnose `--iterations` for repetition.
- `--shuffle-seed` adds random `-shuffle=<seed>` per iteration and records seeds in `report.json`.
- `--parallel-iterations <N>` runs up to N diagnose iterations at once. Each worker gets its own ephemeral Postgres and writes a distinct `iteration-<n>.log.jsonl`. It is rejected with `--database-url`.
- `--fail-fast` stops after any failed iteration.
- `--fail-fast-on=<categories>` stops only when an iteration matches comma-separated categories: `failure`, `timeout`, `slow`, or `any`. Use it when the hypothesis is specific. Do not use `slow` unless `--slow-threshold` is set intentionally.
</semantics>

<defaults>
- Flake hunt: `--iterations 25`, `--timeout 10m` in go test, single package.
- Timeout hunt: `--iterations 5 --fail-fast-on=timeout`, short `--timeout` in go test.
- Slow hunt: `--iterations 3 --fail-fast-on=slow`, `--slow-threshold 5s` on diagnose.
- One-test isolation: `--iterations 100`, `--run '^TestName$' ./path` after `--`.
- Fast narrow rerun: `--iterations 100 --parallel-iterations 4`, single non-intensive test only.
</defaults>

`--ai-output`: first stdout line is the results directory; last stdout line is `<resultsDir>/report.json`. Capture both.

Postgres is ephemeral. Serial diagnose uses one prepared Postgres and restores the prepared snapshot between iterations. Parallel diagnose uses one prepared Postgres per worker and restores that worker's snapshot before reuse. Ignore cross-iteration DB pollution when comparing iterations. Suspect intra-iteration pollution between tests in one package.

Reuse existing DB: `--database-url postgres://...`.
Do not combine `--database-url` with `--parallel-iterations > 1`.
Ctrl+C still analyzes partial results.

<long_running>
For commands expected to exceed 2 minutes:
- Start command with `--ai-output`.
- Use background execution.
- Perform exactly one smoke check of terminal output to catch immediate setup failures.
- Do not poll or wait for completion unless user explicitly asks.
- Tell user the run is continuing in background and analysis will resume when report path appears.
- On completion, read `report.json` first; read raw logs only for entries in `timeouts`, then `flakes`, then relevant `slow`.
</long_running>

<output_layout>
Under repo root.

`diagnose-<targetSlug>-<config>-<YYYYMMDDHHMMSS>/`

- `<targetSlug>`: from trailing package patterns. Leading `./` stripped, `/...` becomes `_allpkgs`, bare `...` becomes `allpkgs`, `/` becomes `_`.
- `<config>`: `it<N>`, optional `p<N>` when `--parallel-iterations > 1`, `h` + 8-hex hash of the full go test argument list, optional `ff`, optional `ffon<categories>` for `--fail-fast-on`, `shuffle`, and optional `slow<duration>` when `--slow-threshold` differs from default. Long basenames shorten slug and drop optional tokens.

Inner layout:

```
diagnose-<targetSlug>-<config>-<YYYYMMDDHHMMSS>/
├── iteration-<n>.log.jsonl
├── report.json
├── report.csv
└── logs/
    └── <short-pkg>_<test>_iter-<n>.log
```
</output_layout>
</diagnose_cli>

<profiles>
When logs and `report.json` are insufficient, narrow with `-run` after `--`, then add standard `go test` profile flags.

<flags>
- `-race`: suspected data race. Heavy. See `<D name="race">`.
- `-cpuprofile`, `-memprofile`: CPU hotspots, allocation pressure, slow tests.
- `-blockprofile`, `-mutexprofile`: channel blocking, lock blocking, mutex contention, hangs.
- `-trace trace.out`: scheduler stalls, long GC pauses, end-to-end timing; open with `go tool trace`.
</flags>

View profiles with `go tool pprof -http=:0 cpu.prof` or profile equivalent. View traces with `go tool trace trace.out`.
Set `GOTRACEBACK=all` or `single` on the test process for fuller panic/fatal stacks.
For stuck tests, send SIGQUIT (`Ctrl+\` on many terminals) to the `go test` process to dump goroutines.

Each `diagnose` iteration runs fresh `go test`. Fixed profile paths overwrite. Use `--iterations 1` while profiling or distinct paths per run. For deep profile work, run a single narrowed `go test`.

Use `GODEBUG=gctrace=1` or `schedtrace` only after isolating a small `-run`; output is noisy.
</profiles>

<reports>
Primary source: `report.json`.

```sh
jq '.' <resultsDir>/report.json
```

<test_entry_fields>
- `package`, `test`
- `runs`, `successes`, `fails`, `skips`, `timeouts`
- `min_elapsed`, `max_elapsed`, `p50_elapsed` (nanoseconds)
- `iterations` (indexes test ran in)
- `logs` (array of problem logs with `type`, `iters`, and `path` pattern relative to resultsDir)
</test_entry_fields>

Top-level buckets: `flakes` (mixed pass/fail), `failures` (always failed), `timeouts` (hit `-timeout`), `slow` (exceeded `--slow-threshold`).
CSV = same data, worst-first, human-skimmable.

<narrow>
When many tests are flagged, pick one before diagnosing.
1. Show user top-N from CSV.
2. Ask which to focus on.
3. Read that test's `logs` paths, replacing `{iter}` with the specific iteration you want to inspect.

```sh
ls <resultsDir>/logs | grep <sanitized-test-name>
cat <resultsDir>/logs/<file>.log
```
</narrow>
</reports>

<diagnose>
Match logs and stats against playbook. State hypothesis before suggesting fix.

<playbook>

<A name="isolate">
Pass alone, fail in package: another test corrupts state. In Chainlink, suspect shared Postgres first. `diagnose` restores prepared DB snapshots between iterations, not between tests in one iteration.

```sh
# Many iterations: serial diagnose restores its prepared DB snapshot between iterations.
# Add --parallel-iterations only for narrow, non-intensive tests.
go -C tools/test run . diagnose --iterations 100 -- --run '^TestName$' ./path/to/package
go -C tools/test run . diagnose --iterations 100 --parallel-iterations 4 -- --run '^TestName$' ./path/to/package
# Use high count for more efficient runtime if DB resets are unnecessary
go -C tools/test run . diagnose -- --run '^TestName$' -count=100 ./path/to/package
# Use -race flag to help induce more unusual timing
go -C tools/test run . diagnose -- --run '^TestName$' -count=100 -race ./path/to/package
```

If still flaky alone, root cause is inside test or code under test.
</A>

<B name="package">
```sh
go -C tools/test run . diagnose --iterations 50 -- ./path/to/package
```
Reproduces in package but not isolation: cross-test dependency.
Common culprits:
- Shared DB rows/tables missing `t.Cleanup` deletion.
- Package-level `var` singletons (keystores, caches, registries).
- Global logger, metric, feature-flag state.
- Shared mock servers without reset.
</B>

<C name="order">
```sh
go test -shuffle=on -count=50 -failfast ./path/to/package
```
Shuffle changes pass rate: order matters. Fix like `<B name="package">`. Capture failing seed (`go test -shuffle=<seed>`) and give it to user.
</C>

<D name="race">
Triggers: stack trace lines do not match `t.Fatal`, nil-pointer panic on unreachable path, inconsistent field values.
```sh
go -C tools/test run . diagnose --iterations 20 -- --race --run '^TestName$' ./path/to/package
```
`-race` is slow and memory-heavy. Use after hypothesis, narrowed with `-run`.
</D>

<E name="resources">
Symptom: fails under load or CI only.
```sh
go test -cpu=1,2,4 -count=20 -failfast ./path/to/package
go test -parallel=1 -count=20 -failfast ./path/to/package
```
Heavy parallelism plus single Postgres can starve connections and cause spurious timeouts.
</E>

<F name="timeout">
For `timeouts` bucket:
- Open `<resultsDir>/logs/<...>_iter-N.log`. `panic: test timed out` includes `running tests:` with active tests at timeout. Analyzer re-attributes; raw log still has goroutine stacks.
- Look for chan receive, `sync.WaitGroup.Wait`, `testutils.WaitTimeout` blocking forever.
- Check service dependencies (Postgres, local server, mock clock) for wrong state.
</F>

<G name="slow">
For `slow` bucket:
- Compare `p50_elapsed` vs `max_elapsed`. Wide spread = intermittent slow (I/O, retries). Narrow spread = test is heavy; rescope.
- Look for `time.Sleep`, long polling loops, retry helpers with generous defaults.
- Chainlink suspects: on-chain event waits with coarse intervals, reconcile loops, long OCR rounds.
</G>

</playbook>
</diagnose>

<fixes>
Lead with hypothesis. Pick one fix archetype:
- Missing cleanup: add `t.Cleanup(func() { ... })` for rows, connections, singletons.
- Global state: move to per-test constructor, or guard and reset in `TestMain`.
- Timing assumption: replace sleeps with `require.eventually`, channel sync.
- Race: narrow shared field, use `sync.Mutex` / `atomic.*`, or redesign sharing.
- DB contention: use separate schema/user per test; package-level `sync.Mutex` on affected tables only as last resort.
- Dead flake on dead code: delete test. See `tools/test/fixing-flaky-tests.md` section 8.

Show contextual diff. Do not describe fix abstractly.

<restricted>
- Do not modify the test's goal to make it pass.
- Do not remove tests or assertions unless replacing them with better ones or deleting dead code with high confidence.
</restricted>
</fixes>

<verify>
Re-run same-scope `diagnose` after fix.
```sh
go -C tools/test run . diagnose --iterations <N> -- <same go test args as before>
```
Compare new `report.json` against previous. Success: target absent from `flakes`, `failures`, `timeouts`, and `slow`. If still present, undo own fix, revise hypothesis, repeat analysis.
</verify>

<chainlink>
- Serial diagnose uses one ephemeral Postgres and restores its prepared snapshot between iterations. Parallel diagnose uses one ephemeral Postgres per worker and restores that worker before reuse. Neither mode resets between tests within one iteration. First suspect for pass-alone / fail-in-package.
- Avoid `--parallel-iterations` for broad `core/...`, `deployment/...`, `--race`, profiling/tracing with fixed output paths, constrained CI workers, DB contention investigations, and `--database-url`.
- `core/internal/testutils` helpers: `testutils.NewTestDB`, `testutils.AssertEventually`, `pgtest.NewSqlxDB`. Prefer over hand-rolled.
- `t.Parallel()` plus one DB can exhaust connections. Pattern: `connection refused` or deadline-exceeded DB calls. Remove `t.Parallel()` from hottest subtests before scaling DB.
- Simulated-chain tests (`backends.NewSimulatedBackend`, `simchain`) are frequent slow offenders. Check `time.Sleep` inside mining loops.
- Default `diagnose` scope is one package or one subtree. Never use `./core/...` without approval.
</chainlink>

<skip>
Do not use this skill when:
- User has known fix. Apply directly. Ask if they want you to verify.
- Test fails deterministically on first run. Use normal debug, not multi-run `diagnose`.
- User wants full-suite CI prep. Use `go -C tools/test run . run` or `go -C tools/test run . gotestsum` (or `make new_test` / `make new_gotestsum`).
</skip>
