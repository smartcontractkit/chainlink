---
name: chainlink-test-diagnosis
description: >-
  Diagnoses and fixes unstable Chainlink Go tests (flakes, races, timeouts, deadlocks,
  slow runs). Use for non-deterministic failures, CI-only instability, or test runtime.
  Do NOT use for deterministic failures, routine runs, or full-suite CI prep.
---

<absolute_constraints>
- DO NOT use this skill if the user already has a known fix (apply it directly).
- DO NOT use for deterministic first-run failures (use normal debug).
- DO NOT use for full-suite CI prep (use `make new_test` or `make new_gotestsum` instead).
- ONLY run tests in these packages without explicit user approval: `core/`, `deployment/`. Warn the user if running outside these.
- DO NOT modify the test's core goal to make it pass.
- DO NOT remove tests/assertions unless replacing with better ones or deleting confirmed dead code.
- DO NOT modify package-wide helpers (`testutils`) to fix localized tests.
- IF Postgres sandbox error occurs (`operation not permitted`), ask the user to run the command or approve unsandboxed execution.
- For runs expected >2m: Execute in background. Perform a single 30s crash check, then suspend task and wait for the report.json system notification. DO NOT poll.
</absolute_constraints>

## Purpose
Root-cause analysis and fixes for unstable or slow Chainlink Go tests.
Flow: Scope target -> Check for recent run -> Run `diagnose` -> Analyze -> Apply fix -> Verify.

## Preflight
- Ask for test, package, issue, or permission to discover.
- Start bounded: single test, package, or subtree; use `--fail-fast` or low `--iterations`.
- Classify hypothesis: flake, timeout, slow, panic, deadlock, or race.

<cli_reference>
Authoritative behavior: `go -C tools/test run . diagnose -h`. 
Run from repo root. Flags before `--` belong to the harness. Flags/packages after `--` belong to `go test`.
Always use `--ai-output`.

**Commands:**
- Help: `go -C tools/test run . diagnose -h`
- Flake hunt (default): `go -C tools/test run . diagnose --iterations 25 -- --timeout 10m ./path/to/package`
- Timeout hunt: `go -C tools/test run . diagnose --iterations 5 --fail-fast-on=timeout -- --timeout 2m ./path/to/package`
- Slow hunt: `go -C tools/test run . diagnose --iterations 3 --fail-fast-on=slow --slow-threshold 5s -- ./path/to/package`
- One-test isolation: `go -C tools/test run . diagnose --iterations 100 -- --run '^TestName$' ./path`
- Fast narrow rerun: `go -C tools/test run . diagnose --iterations 100 --parallel-iterations 4 -- --run '^TestName$' ./path`
- Shuffle test order: `go test -shuffle=on -count=50 -failfast ./path/to/package`
- Detect Race: `go -C tools/test run . diagnose --iterations 20 -- --race --run '^TestName$' ./path/to/package`
- CPU/Memory load: `go test -cpu=1,2,4 -count=20 -failfast ./path/to/package`
- Verify fix: `go -C tools/test run . diagnose --iterations <N> -- <same go test args>`
- Lint check: `golangci-lint run ./<packages-you-change> --fix`
</cli_reference>

## Execution & Analysis
- **Postgres:** Serial diagnose restores DB between iterations. Parallel gives each worker an ephemeral DB. Neither resets between tests *within* one iteration.
- **Report Analysis:** Read `<resultsDir>/report.json` using `jq`. Top-level buckets: `flakes`, `failures`, `timeouts`, `slow`.
- **Narrowing:** If many tests flag, look for similarities in their failures. If found, present that to the user and ask if they want to continue with that assumption. If not, try to focus on the most problematic test.
- **Profiles:** When logs/report are insufficient, use standard `go test` profile flags (`-race`, `-cpuprofile`, `-trace`, etc.). View with `go tool pprof` or `go tool trace`.

<logs>

### `diagnose` outputs

```
diagnose-results/
|-- iteration-n.log.jsonl # DO NOT READ unless absolutely necessary; full log outputs, long and messy
|-- postgres-state-n.md # Final state of postgres DB after test iteration. Read if diagnosing DB-based errors or hangs.
|-- report.json # Read this; summary of full `diagnose` run
|-- report.csv # DO NOT READ; human readable csv
|-- logs/ # Extracted individual test logs
|---- pkg_TestName_iter-n.log # Logs for individual slow/failing test
```

Spawn a sub-agent to read specific log files from the end up.
You MUST output ONLY valid JSON matching this exact structure. No markdown formatting, no explanations, no yapping:
{
  "logs_read": ["log_path_1.log", "log_path_2.log"],
  "failure_diagnosis": [
    {
      "possible_reason": "explanation",
      "evidence": "reasoning and evidence"
    }
  ]
}
</logs>

## Playbook & Fixes
Lead with your hypothesis before writing code. Show contextual diffs, do not describe fixes abstractly.

1. **Isolate (Pass alone, fail in package):** Cross-test dependency. Missing `t.Cleanup`, global state (`var` singletons, loggers), or shared mock servers. Fix by moving state to per-test constructors or using `t.Cleanup`.
2. **Order (Shuffle changes pass rate):** Same as isolation. Fix cross-test leakage. Capture failing seed and provide to user.
3. **Race:** Triggers on weird stack traces or nil pointers. Use `-race`. Fix with `sync.Mutex`, `atomic.*`, or narrow shared fields.
4. **Timeout:** Check logs for blocking (chan receive, `Wait`, `testutils.WaitTimeout`). Use `synctest` to improve tests relying on channels.
5. **Slow:** Compare `p50` vs `max_elapsed`. Look for `time.Sleep` or coarse polling loops. Replace with `require.eventually` or channel sync. Simulated chains are frequent offenders.
6. **Resources:** If failing under load/CI only, DB connections might be exhausted by `t.Parallel()`. Use separate schema/user per test.

<context_compaction>
When summarizing context, strictly maintain state in this format:

## [TestName]
Failure: [suspected failure reasons]
SuspectedFix: [the fix you've implemented or want to try]
NextStep: [the next step for diagnosing/fixing/verifying the test]
</context_compaction>