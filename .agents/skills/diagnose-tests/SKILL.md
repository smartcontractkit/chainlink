---
name: diagnose-and-fix-flaky-slow-tests
description: >-
  Repeatedly runs Chainlink Go unit tests with the `survey` command, parses the
  flake/failure/timeout/slow report it emits, and helps the user root-cause and
  fix a specific flaky or slow test.
---

<purpose>
Root-cause flake, timeout, or slow test. Three phases: survey (gather evidence), analyze (read reports), diagnose (apply playbook, propose fix).
</purpose>

<before>
- Ask user if they have a specific test, package, or issue they're working on, or if they want to just explore and discover new issues.
- `./core/...` survey = 5m-10m per iteration. Never survey unbounded target without explicit approval. Start with `--fail-fast` or small iterations.
- Confirm hypothesis: flake, timeout, slow, or panic. Decides flags.
</before>

<survey>
<restrictions>
If you are in sandbox mode, you will be unable to run the `survey` command. Either prompt the user to run it themselves, or exit you from sandbox mode. If you hit the following error, it's likely you're restricted by sandbox mode.
<error>
Failed to reset database:unable to drop postgres database: failed to connect to `host=localhost user=postgres database=template1`: dial error (dial tcp 127.0.0.1:55001: connect: operation not permitted)
</error>
</restrictions>

Run from repo root:

```sh
# Command help
go -C ./tools/test run . survey -h
# Example command
go -C ./tools/test run . survey --iterations <N> --slow-threshold <duration> --timeout <duration> --run '<regex>' --race --ai-output <target-pattern>
```

Flag semantics:
- `--run <regex>` → passed as `go test -run=<regex>`. Narrow to a specific test or subtest. Use `^TestName$` for exact match, `TestName/sub` for a subtest, or `TestA|TestB` for union.
- `--race` → passed as `go test -race`. Use after isolating a hypothesis; slow and memory-heavy.
- Defaults by hypothesis:
  - Flake hunt: `--iterations 25 --timeout 10m`, single package.
  - Timeout hunt: `--iterations 5 --timeout 2m`. Short per-iter timeout surfaces hang fast.
  - Slow hunt: `--iterations 3 --slow-threshold 5s`.
  - One-test isolation repro: `--iterations 100 --run '^TestName$' ./path/to/package`.

`--ai-output` prints results directory path to stdout. Capture it.

Postgres is shared, ephemeral. Survey snapshots after setup, restores before each iteration. Cross-iteration DB pollution: not a concern. Intra-iteration pollution between tests in same package: common cause — see playbook.

Reuse existing DB: `--database-url postgres://…`.

Ctrl+C still runs analysis on partial results.

Output layout (under repo root):

Directory basename:

`survey-<targetSlug>-<config>-<YYYYMMDDHHMMSS>/`

- **`<targetSlug>`** — From the go test package pattern: leading `./` stripped, `/...` becomes `_allpkgs`, bare `...` becomes `allpkgs`, `/` becomes `_`, other awkward characters become `_`.
- **`<config>`** — Hyphen-separated survey flags, e.g. `it<N>`, `to<duration>`, optional `race`, `shuffle`, `ff`, `p<N>`, `cpu-…`, optional `r<sanitized-run>-<8hex>` when `--run` is set, optional `slow<duration>` when `--slow-threshold` differs from the default. If the full basename would exceed ~220 bytes, the slug is shortened and optional tokens are dropped in phases (the date suffix and `-run` hash stay for disambiguation).

Inner layout:

```
survey-<targetSlug>-<config>-<YYYYMMDDHHMMSS>/
├── iteration-<n>.log.jsonl
├── report.json
├── report.csv
└── logs/
    └── <short-pkg>_<test>_iter-<n>.log
```
</survey>

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
Pass alone, fail in package: other test corrupts state. Chainlink: usually shared Postgres (survey restores between iterations, not between tests in one iteration).

```sh
go -C ./tools/test run . survey --iterations 100 --run '^TestName$' ./path/to/package
```

Still flakes alone: problem inside test or code under test.
</A>

<B name="package">
```sh
go -C ./tools/test run . survey --iterations 50 ./path/to/package
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
go -C ./tools/test run . survey --iterations 20 --race --run '^TestName$' ./path/to/package
```
`-race` costly (slow + memory-heavy). Use after hypothesis, narrowed with `--run`.
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
Re-run same-scope survey after fix:
```sh
go -C ./tools/test run . survey --iterations <N> <target>
```
Compare new `report.json` vs previous. Success: test absent from `flakes`, `failures`, `timeouts`, `slow`. Still present → revert, revise hypothesis, repeat diagnose.
</verify>

<chainlink>
- Single shared Postgres across `core/...`. Survey restores between iterations, not between tests within iteration. First suspect for pass-alone / fail-in-package.
- `core/internal/testutils` helpers: `testutils.NewTestDB`, `testutils.AssertEventually`, `pgtest.NewSqlxDB`. Prefer over hand-rolled.
- `t.Parallel()` + one DB: exhausts connections. Flake pattern `connection refused` or deadline-exceeded in DB calls → remove `t.Parallel()` from hottest subtests rather than scale DB.
- Simulated-chain tests (`backends.NewSimulatedBackend`, `simchain`) = frequent slow offenders. Check `time.Sleep` inside mining loops.
- Default survey scope = one package or one subtree. Never `./core/...` without approval.
</chainlink>

<skip>
Do not use this skill when:
- User has known fix — apply directly.
- Test fails deterministically first run — normal debug, no survey loop.
- User wants full-suite CI prep — use `test` or `gotestsum` subcommands.
</skip>
