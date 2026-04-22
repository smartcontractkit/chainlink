---
name: diagnose-and-fix-flaky-slow-tests
description: >-
  Repeatedly runs Chainlink Go unit tests with the `diagnose` command, parses the
  flake/failure/timeout/slow report it emits, and helps the user root-cause and
  fix a specific flaky or slow test.
---

<purpose>
Root-cause flake, timeout, or slow test. Three phases: run `diagnose` to gather evidence, analyze reports, then apply the playbook and fix.
</purpose>

<before>
- Ask user if they have a specific test, package, or issue they're working on, or if they want to just explore and discover new issues.
- `./core/...` with `diagnose` = 5m-10m per iteration. Never run `diagnose` on an unbounded target without explicit approval. Start with `--fail-fast` or small iterations.
- Confirm hypothesis: flake, timeout, slow, or panic. Decides `go test` flags after `--`.
</before>

<diagnose_cli>
<restrictions>
If you are in sandbox mode, you will be unable to run the `diagnose` command. Either prompt the user to run it themselves, or exit you from sandbox mode. If you hit the following error, it's likely you're restricted by sandbox mode.
<error>
Failed to reset database:unable to drop postgres database: failed to connect to `host=localhost user=postgres database=template1`: dial error (dial tcp 127.0.0.1:55001: connect: operation not permitted)
</error>
</restrictions>

Run from repo root. **Harness-only flags** (before `--`): `--iterations`, `--slow-threshold`, `--fail-fast`, `--shuffle-seed`. **Everything after `--`** is passed to `go test` (e.g. `-timeout`, `-race`, `-run`, package patterns). Put **package patterns last** (usual `go test` layout).

```sh
# Command help
go -C ./tools/test run . diagnose -h
# Example: harness flags, then --, then go test flags and packages
go -C ./tools/test run . diagnose --iterations <N> --slow-threshold <duration> --fail-fast --ai-output -- --timeout <duration> --run '<regex>' --race ./path/to/package/...
```

Harness semantics:
- Prepends `go test -json` each iteration (drops duplicate `-json`); adds `-count=1` unless you set `-count` greater than 1 on `go test` (prefer diagnose `--iterations` for repetition).
- `--shuffle-seed` → random `-shuffle=<seed>` per iteration, recorded in `report.json`.
- Defaults by hypothesis (tune `go test` after `--`):
  - Flake hunt: `--iterations 25`, `--timeout 10m` in go test, single package.
  - Timeout hunt: `--iterations 5`, short `--timeout` in go test.
  - Slow hunt: `--iterations 3`, `--slow-threshold 5s` on diagnose.
  - One-test isolation: `--iterations 100`, `--run '^TestName$' ./path` after `--`.

`--ai-output` prints results directory path to stdout. Capture it.

Postgres is shared, ephemeral. Diagnose snapshots after setup, restores before each iteration. Cross-iteration DB pollution: not a concern. Intra-iteration pollution between tests in same package: common cause — see playbook.

Reuse existing DB: `--database-url postgres://…`.

Ctrl+C still runs analysis on partial results.

Output layout (under repo root):

Directory basename:

`diagnose-<targetSlug>-<config>-<YYYYMMDDHHMMSS>/`

- **`<targetSlug>`** — From trailing package patterns (put packages last). Leading `./` stripped, `/...` becomes `_allpkgs`, bare `...` becomes `allpkgs`, `/` becomes `_`.
- **`<config>`** — `it<N>`, `h` + 8-hex hash of the full go test argument list (disambiguates flags), optional `ff`, `shuffle`, optional `slow<duration>` when `--slow-threshold` differs from the default. If the full basename would exceed ~220 bytes, the slug is shortened and optional tokens are dropped in phases.

Inner layout:

```
diagnose-<targetSlug>-<config>-<YYYYMMDDHHMMSS>/
├── iteration-<n>.log.jsonl
├── report.json
├── report.csv
└── logs/
    └── <short-pkg>_<test>_iter-<n>.log
```
</diagnose_cli>

<profiles>
When logs and `report.json` are not enough, suggest **narrowing** with `-run` (after `--`), then adding **standard `go test` profile flags** so the user gets artifacts to inspect offline.

**Typical flags** (all go in the `go test` argument list after `--`):

| Flag | Use when |
|------|----------|
| `-race` | Suspected data race (see playbook `<D name="race">`); heavy. |
| `-cpuprofile`, `-memprofile` | CPU hotspots, alloc pressure, slow tests. |
| `-blockprofile`, `-mutexprofile` | Blocking on channels/locks, mutex contention, “hangy” tests. |
| `-trace trace.out` | Scheduler stalls, long GC pauses, end-to-end timing; open with `go tool trace`. |

**Viewing**: `go tool pprof -http=:0 cpu.prof` (or `mem.prof`, `block.prof`, `mutex.prof`). For `-trace`, `go tool trace trace.out`.

**Environment**: `GOTRACEBACK=all` (or `single`) on the test process for fuller stacks on panic/fatal. On a **stuck** test, sending **SIGQUIT** (`Ctrl+\` on many terminals) to the `go test` process dumps all goroutines—useful for deadlocks.

**`diagnose` + profiles**: Each iteration runs a fresh `go test`. If you pass a fixed `-cpuprofile=cpu.prof`, later iterations **overwrite** the same file. Prefer **`--iterations 1`** while profiling, or use **distinct paths per run** (shell loop / manual re-runs). For deep profile analysis, a single `go test` invocation with `-run` is often clearer than a long multi-iteration `diagnose`.

**Noise**: `GODEBUG=gctrace=1`, `schedtrace`—very verbose; only after isolating a small `-run`.
</profiles>

<reports>
Primary source: `report.json`.

```sh
jq '.' <resultsDir>/report.json
```

TestEntry fields:
- `package`, `test`
- `runs`, `successes`, `fails`, `skips`, `timeouts`
- `min_elapsed`, `max_elapsed`, `p50_elapsed` (nanoseconds)
- `iterations` (indexes test ran in)
- `log_files` (paths relative to resultsDir)

Top-level buckets: `flakes` (mixed pass/fail), `failures` (always failed), `timeouts` (hit `-timeout`), `slow` (exceeded `--slow-threshold`).

CSV = same data, worst-first, human-skimmable.

<narrow>
Many flagged tests expected. Pick one before diagnosing.
1. Show user top-N from CSV.
2. Ask which to focus on.
3. Read that test's `log_files`.

```sh
ls <resultsDir>/logs | grep <sanitized-test-name>
cat <resultsDir>/logs/<file>.log
```
</narrow>
</reports>

<diagnose>
Pattern-match logs + stats against playbook. State explicit hypothesis before suggesting fix.

<playbook>

<A name="isolate">
Pass alone, fail in package: other test corrupts state. Chainlink: usually shared Postgres (`diagnose` restores between iterations, not between tests in one iteration).

```sh
go -C ./tools/test run . diagnose --iterations 100 -- --run '^TestName$' ./path/to/package
```

Still flakes alone: problem inside test or code under test.
</A>

<B name="package">
```sh
go -C ./tools/test run . diagnose --iterations 50 -- ./path/to/package
```
Reproduces here but not isolation: cross-test dependency. Common chainlink culprits:
- Shared DB rows/tables missing `t.Cleanup` deletion.
- Package-level `var` singletons (keystores, caches, registries).
- Global logger, metric, feature-flag state.
- Shared mock servers without reset.
</B>

<C name="order">
```sh
go test -shuffle=on -count=50 -failfast ./path/to/package
```
Shuffle changes pass rate: order matters. Fixes = §B. Capture seed from failing run (`go test -shuffle=<seed>`). Give seed to user.
</C>

<D name="race">
Trigger: stack trace lines don't match `t.Fatal`; nil-pointer panic on unreachable path; inconsistent field values.
```sh
go -C ./tools/test run . diagnose --iterations 20 -- --race --run '^TestName$' ./path/to/package
```
`-race` costly (slow + memory-heavy). Use after hypothesis, narrowed with `-run`.
</D>

<E name="resources">
Symptom: fails under load or in CI only.
```sh
go test -cpu=1,2,4 -count=20 -failfast ./path/to/package
go test -parallel=1 -count=20 -failfast ./path/to/package
```
Heavy parallelism + single Postgres → connection starvation → spurious timeouts.
</E>

<F name="timeout">
For `timeouts` bucket:
- Open `<resultsDir>/logs/<...>_iter-N.log`. `panic: test timed out` block has `running tests:` section listing active tests at timeout. Analyzer already re-attributes; raw log still has goroutine stacks.
- Look for chan receive, `sync.WaitGroup.Wait`, `testutils.WaitTimeout` blocking forever.
- Check service dependencies (Postgres, local server, mock clock) for wrong state.
</F>

<G name="slow">
For `slow` bucket:
- Compare `p50_elapsed` vs `max_elapsed`. Wide spread = intermittent slow (I/O, retries). Narrow spread = test is heavy; rescope.
- Look for `time.Sleep`, long polling loops, retry helpers with generous defaults.
- Chainlink: on-chain event waits with coarse intervals, reconcile loops, long OCR rounds.
</G>

</playbook>
</diagnose>

<fixes>
Lead with hypothesis. Pick one archetype:
- Cleanup missing → add `t.Cleanup(func() { … })` for rows, connections, singletons.
- Global state → move to per-test constructor, or guard + reset in `TestMain`.
- Timing assumption → replace sleeps with `gomega.Eventually`, `testutils.AssertEventually`, channel sync.
- Race → narrow shared field, `sync.Mutex` / `atomic.*`, or redesign sharing.
- DB contention → separate schema/user per test; package-level `sync.Mutex` on affected tables as last resort.
- Dead flake on dead code → delete test. See `tools/test/fixing-flaky-tests.md` §8.

Show diff in context (Read → Edit). Do not describe fix abstractly.
</fixes>

<verify>
Re-run the same-scope `diagnose` run after fix:
```sh
go -C ./tools/test run . diagnose --iterations <N> -- <same go test args as before>
```
Compare new `report.json` vs previous. Success: test absent from `flakes`, `failures`, `timeouts`, `slow`. Still present → revert, revise hypothesis, repeat root-cause analysis.
</verify>

<chainlink>
- Single shared Postgres across `core/...`. `diagnose` restores between iterations, not between tests within iteration. First suspect for pass-alone / fail-in-package.
- `core/internal/testutils` helpers: `testutils.NewTestDB`, `testutils.AssertEventually`, `pgtest.NewSqlxDB`. Prefer over hand-rolled.
- `t.Parallel()` + one DB: exhausts connections. Flake pattern `connection refused` or deadline-exceeded in DB calls → remove `t.Parallel()` from hottest subtests rather than scale DB.
- Simulated-chain tests (`backends.NewSimulatedBackend`, `simchain`) = frequent slow offenders. Check `time.Sleep` inside mining loops.
- Default `diagnose` scope = one package or one subtree. Never `./core/...` without approval.
</chainlink>

<skip>
Do not use this skill when:
- User has known fix — apply directly.
- Test fails deterministically first run — normal debug, no multi-run `diagnose` loop.
- User wants full-suite CI prep — use `test` or `gotestsum` subcommands.
</skip>
