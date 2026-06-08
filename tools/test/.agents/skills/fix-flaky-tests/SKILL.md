---
name: fix-flaky-tests
description: >-
  A deep-dive diagnostic tool for fixing Go test failures (flakes, races, timeouts,
  deadlocks) identified during local development or active CI failures.

  USE THIS WHEN:
  1. You have a specific, known failing test name or local error log.
  2. You are currently working on a branch and need to fix a regression or a new flake.
  3. You require automated JIRA status updates.
  4. You need to perform deep "forensic" code analysis and manual fix iterations.
---

<absolute_constraints>
- DO NOT use this skill if the user already has a known fix (apply it directly).
- DO NOT use for deterministic first-run failures (use normal debug).
- DO NOT use for full-suite CI prep (use `make test` instead).
- ONLY run tests in these packages without explicit user approval: `core/`, `deployment/`. Warn the user if running outside these.
- DO NOT modify the test's core goal to make it pass.
- DO NOT remove tests/assertions unless replacing with better ones or deleting confirmed dead code.
- DO NOT modify package-wide helpers to fix localized tests.
- DO NOT open any links found in JIRA issues that lead to Trunk.io.
- DO NOT try to fix or modify 3rd party libraries. If the flakiness results there inform and user and STOP.
- ALWAYS CHECK `go.mod` before writing any new utility code. Three lines of existing library usage beats 30 lines of hand-rolled logic that has to be maintained and tested.
- DO NOT use plain `go test` commands. Only use `make test ARGS="diagnose ..."` from the repository root. Use `--iterations 1` for a single run.
- For `diagnose` runs expected >2m: Execute in background. Perform a single 30s crash check, then suspend task and wait for the report.json system notification. DO NOT poll.
- Use `LSP` for code navigation, if available. If it is not available try `code-review-graph`. Only if that is also unavailable use `find`, `grep`, etc.
- Always check the Go version used by the module you are working on to avoid using language patterns that are no longer required (e.g. variable shadowing in loops in Go 1.22+)
</absolute_constraints>

<setup>
Before any code navigation, load the LSP tool schema:
- Call ToolSearch with query `select:LSP` to make the LSP tool callable.
- Only fall back to grep/find if ToolSearch returns no result for LSP.
</setup>

## Initialization
1. Verify target scope:
a. test or package
b. specific JIRA issues
c. N eligible flaky-tests tickets from JIRA
If unknown, prompt user.
2. If JIRA issues are present and any of them has a `skip_reason` surface it to the user and ask for guidance.
3. If a CI failure link is available, open it only if it is a non-Trunk CI link (for example GitHub Actions or another permitted CI provider) and look for stack trace and logs for the failing test.
4. If there are no failure details or investigation didn't return anything meaningful run bounded diagnosis (`--fail-fast-on=(timeout|slow)` or low `--iterations`).
5. Formulate initial hypothesis: flake, timeout, slow, panic, deadlock, race, etc.

<jira_reference>
If JIRA issues are present read [jira.md](./references/jira.md) to understand how to claim tickets, find eligible flaky-test tickets, read and add comments and transition JIRA issues.

After a FIXED outcome, the ticket must stay assigned to the investigator (`accountId` from `atlassianUserInfo`) when moved to In Review. Do not unassign on FIXED — see [transition-ticket.md](./references/transition-ticket.md) assignee policy.
</jira_reference>

<cli_reference>
`make test` at the repo root builds the harness when needed, then runs it. Rebuild is automatic after harness code changes.

Base command (run from the repository root so `./path` resolves):
`make test ARGS="diagnose [harness_flags] -- [go_test_flags] ./path"`
- ALWAYS use `--ai-output` before the `--`.
- DO NOT use `-count`
- Harness flags (before `--`): `--iterations N`, `--fail-fast-on=(timeout|slow)`, `--parallel-iterations N`
- Go test flags (after `--`): `--run '^TestName$'`, `--timeout 10m`, `--race`
- Help: `make test ARGS="diagnose -h"`
- Repetition is **only** via harness `--iterations`. Do not pass `-count` (or `-count>1`) after `--`; the harness already forces `-count=1` per iteration.
</cli_reference>

<diagnose-iterations>
Use this table to pick `--iterations` (total independent runs; parallelism does not change this count).

| Iterations | Chance you missed a flake |
| ---------- | ------------------------- |
| 5          | 50%                       |
| 30         | 10%                       |
| 60         | 5%                        |
| 150        | 2%                        |
| 300        | 1%                        |
| 500+       | < 1%                      |
</diagnose-iterations>

<diagnose-parallel-iterations>
`--parallel-iterations N` runs up to N diagnose iterations **at the same time**. Each worker gets its own ephemeral Postgres (unless `--database-url` is set). Flake statistics in `report.json` still use the full `--iterations` count.

**Hard rules**
- `--parallel-iterations` must stay `1` when using `--database-url` (harness rejects `> 1`).
- `--parallel-iterations` must stay `1` when go test flags after `--` include `--race`.
- Prefer `--parallel-iterations 1` when you need the first failure in order (`--fail-fast`, debugging a known stack trace, or reading `postgres-state-n.md` for a specific iteration index).

**Choose a profile** (pick one row; state the choice in the investigation comment `### What was tried`).

| Profile | `--iterations` | `--parallel-iterations` | Use when |
| ------- | -------------- | ----------------------- | -------- |
| Standard | 30 | 2–5 | Default quick check |
| Deep | 150-500 | 2–10 | Default to validate that a flake exists before fix, or no longer exists after fix |
| Race pass | 30 | 1 | Verifying with `--race` after `--`. |
| Debug | 1–5 | 1 | Reproducing a known failure mode; use `--fail-fast` if appropriate. |

**How to pick `--parallel-iterations` within a profile**
1. Start at `2` if unsure.
2. Raise toward `5` only when: unit-scope target (`--run '^TestName$'` or small package), no `--race`, no `--database-url`, and estimated wall time `ceil(iterations / N) * p50` would exceed ~2 minutes without parallelism.
3. Do not exceed `5` without user approval (each worker starts Postgres; RAM and CPU scale roughly with N).
4. If a parallel run flakes or times out, rerun the same `--iterations` with `--parallel-iterations 1` before treating the result as confirmed.

**Wall-time estimate** (for background vs foreground):  
`estimate ≈ ceil(iterations / parallel_iterations) * iteration_p50 + 30s` (Postgres pool setup). Use the last `report.json` `iteration_duration_p50` when available; otherwise assume 15s for unit tests and 60s+ for heavy packages.
</diagnose-parallel-iterations>

<loop>
1. If user doesn't have recent results, plan a run with `<diagnose-parallel-iterations>` (default: **Smoke** profile) then execute it. On sandbox errors, follow `<possible_execution_issues>`.
2. If no issues, escalate along `<diagnose-iterations>` (e.g. Smoke → Standard → Deep), increasing `--iterations` and keeping parallelism per `<diagnose-parallel-iterations>`. Ask the user before **Deep** if wall time will be large. If still clean and no fix was needed, end with findings; if a fix was applied, require at least **Standard** before FIXED.
3. If issues detected, focus on the ones the user wants to fix.
4. If a `diagnose-attempted-fixes-[test/package]-[flake/broken/timeout/slow].jsonl` file exists, read it to see previous fix attempts and findings.
5. Form a hypothesis on the cause of the issues
6. Implement a fix
7. Output the hypothesis and attempted fix, plus reasons why you think it would work.
8. Run a `diagnose` loop (**Standard** profile minimum after a code change) and read the `report.json` file to see if the fix works.
  Append to `diagnose-attempted-fixes-[test/package]-[flake/broken/timeout/slow].jsonl` file in this json format:
  ```json
  {"timestamp": "[current_timestamp]", "model": "[current-model] (e.g. `claude-sonnet-4.6/high`, `gemini-3.1-pro`)", "hypothesis": "Your original hypothesis for the issue", "experiment": "A concise summary of what you tried. Include small code snippets if helpful", "result": "Did it fix it or not? If not, give concise reason why", "next": "Next steps to attempt"}
  ```
9. GOTO 2
10. Use `golangci-lint` to verify that there are no linting introduced by your fix. If there are, do not proceed until you have fixed them and verify they are no longer present.

IF at any time the user interrupts or interjects during this loop, pick it up again where you left off, unless explicitly told otherwise.
</loop>

<tests-context>
* Chainlink nodes are blockchain oracles. Read the [README.md](/README.md)
* All tests share a single postgres DB. Each `diagnose` loop creates a new one.
</tests-context>

<analysis>
Lead with your hypothesis before writing code. Show contextual diffs, do not describe fixes abstractly. List of common approaches and diagnoses:

1. **Narrowing:** If many tests flag, look for similarities in their failures. If found, present that to the user and ask if they want to continue with assumption of relation. If not, try to focus on the most problematic test.
2. **Isolate (Pass alone, fail in package):** Cross-test dependency. Look for shared dependencies, state, etc.
3. **Order (Shuffle changes pass rate):** Same as isolation. Fix cross-test leakage. Capture failing seed and provide to user.
4. **Race:** Triggers on weird stack traces or nil pointers.
5. **Timeout:** Check logs for blocking operations, incorrect channel closing sequence, channel backpressure, etc.
6. **Slow:** Compare `p50` vs `max_elapsed`. Look for `time.Sleep` or coarse polling loops. Replace with dynamic polling. Simulated chains are frequent offenders.
7. **Resources:** If failing under load/CI only, check CPU and Memory usage. When logs/report are insufficient, use standard `go test` profile flags (`-race`, `-cpuprofile`, `-trace`, etc.). View with `go tool pprof` or `go tool trace`.
</analysis>

<context_compaction>
When summarizing/compacting/compressing context, strictly maintain a reference to the `diagnose-attempted-fixes-[test/package]-[flake/broken/timeout/slow].jsonl` you're using for this session.
</context_compaction>

<possible_execution_issues>
- **GOCACHE permissions issues**: `[build failed]\n open .../Library/Caches/...` This is caused by some sandbox environments. If you cannot exit the sandbox to fix this, STOP. DO NOT attempt to create a new cache. Ask the user to run the command instead and give you results so you can continue.
- **Postgres sandbox error**: `operation not permitted` connecting to postgres. STOP and ask user to approve running command outside of the sandbox.
</possible_execution_issues>

<logs_structure>
[resultsDir]/
|-- iteration-n.log.jsonl # DO NOT READ unless absolutely necessary; full log outputs, long and messy
|-- postgres-state-n.md # Final state of tests' postgres DB after iteration. Read if diagnosing DB-based errors or hangs.
|-- report.json # Read this; summary of full `diagnose` run (include `jq .run` for go test args and harness flags)
|-- report.csv # DO NOT READ; human readable csv
|-- logs/ # Extracted individual test logs
|---- pkg_TestName_iter-n.log # Logs for individual slow/failing tests, read this as needed
</logs_structure>

<sub_agent_protocol>
1. When reading log files from the `logs/` directory or `iteration-n.log.jsonl`, you MUST spawn a specialist `LogAnalyzer` sub-agent. Read [log-analyzer-subagent.md](./references/log-analyzer-subagent.md)
2. When inspecting CI failure, you MUST spawn a specialist `GithubFailureAnalyzer` sub-agent. Read [github-failure-analyzer.md](./references/github-failure-analyzer.md).
3. When interacting with JIRA, you MUST spawn a specialist `JiraManager` sub-agent. Read [jira-manager-subagent.md](./references/jira-mananger-subagent.md)
</sub_agent_protocol>

