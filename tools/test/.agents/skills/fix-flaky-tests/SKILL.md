---
name: fix-flaky-tests
description: >-
  Diagnostic tool for fixing Go test failures (flakes, races, timeouts, deadlocks) during local dev or CI.
---

<absolute_constraints>
- DO NOT use this skill if the user already has a known fix (apply it directly).
- IF user wants to only speed up slow test, go directly to [speed-up-tests](./references/speed-up-tests.md).
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
- Use `LSP` for code navigation, if available. Check if it works using a go file from the project. If it is not available try `code-review-graph`. Only if that is also unavailable use `find`, `grep`, etc.
- Always check the Go version used by the module you are working on to avoid using language patterns that are no longer required (e.g. variable shadowing in loops in Go 1.22+)
</absolute_constraints>

<setup>
Call ToolSearch with query `select:LSP` to load LSP tool schema.
Fallback to rg/grep/find only if ToolSearch fails.
</setup>

<initialization>
1. Verify target scope (If unknown, prompt user):
 a. test or package
 b. specific JIRA issues
 c. N eligible flaky-tests tickets from JIRA
2. Before proceeding always ask the user whether she thinks the flake is relatively simple and self-contained or whether it is a complex one that requires a lot of critical thinking and in-depth understanding of the application (e.g. system tests and some integration tests). If it is the latter activate the [complex-investigation-protocol](./references/complex-investigation-protocol.md), before formulating any hypothesis.
3. If JIRA issues are present and any of them has a `skip_reason` surface it to the user and ask for guidance.
4. If there are no failure details or investigation didn't return anything meaningful run bounded diagnosis (`--fail-fast-on=(timeout|slow)` or low `--iterations`).
5. Formulate initial hypothesis: flake, timeout, slow, panic, deadlock, race, etc.
</initialization>

<jira_reference>
Read [jira.md](./references/jira.md) to understand how to claim tickets, find eligible flaky-test tickets, check if there are any tickets related to a specific test, read and add comments and transition JIRA issues.

After a FIXED outcome, the ticket must stay assigned to the investigator (`accountId` from `atlassianUserInfo`) when moved to In Review. Do not unassign on FIXED — see [transition-ticket.md](./references/transition-ticket.md) assignee policy.
</jira_reference>

<cli_reference>
Execute from repository root.
`make test ARGS="diagnose [harness_flags] -- [go_test_flags] ./path"`

- Require `--ai-output` before `--`.
- Forbid `-count`.
- Harness flags: `--iterations N`, `--fail-fast-on=(timeout|slow)`, `--parallel-iterations N`.
- Go test flags: `--run '^TestName$'`, `--timeout 10m`, `--race`.
- Help: `make test ARGS="diagnose -h"`.
- Repetition strictly via `--iterations`.

<diagnose-iterations>
Use iterations for run count. Parallelism does not alter total.
- 5: 50% missed flake
- 30: 10% missed flake
- 60: 5% missed flake
- 150: 2% missed flake
- 300: 1% missed flake
- 500+: <1% missed flake

| Profile   | `--iterations` | `--parallel-iterations` | Use when                                                                          |
| --------- | -------------- | ----------------------- | --------------------------------------------------------------------------------- |
| Quick     | 1-5            | 1-5                     | Quick check to validate no failures                                               |
| Standard  | 30             | 1–5                     | Default standard check                                                            |
| Deep      | 150-500        | 2–10                    | Default to validate that a flake exists before fix, or no longer exists after fix |
| Race pass | 30             | 1                       | Verifying with `--race` after `--`.                                               |
| Debug     | 1–5            | 1                       | Reproducing a known failure mode; use `--fail-fast` if appropriate.               |
</diagnose-iterations>
</cli_reference>

<loop>
1. If user doesn't have recent results, plan a run with `<diagnose-parallel-iterations>` (default: **Smoke** profile) then execute it. On sandbox errors, follow `<possible_execution_issues>`.
2. If no issues, escalate along `<diagnose-iterations>` (e.g. Smoke → Standard → Deep), increasing `--iterations` and keeping parallelism per `<diagnose-parallel-iterations>`. Ask the user before **Deep** if wall time will be large. If still clean and no fix was needed, end with findings; if a fix was applied, require at least **Standard** before FIXED.
3. If issues detected, focus on the ones the user wants to fix.
4. If a `diagnose-attempted-fixes-[test/package]-[flake/broken/timeout/slow].jsonl` file exists, read it to see previous fix attempts and findings.
5. If it is a complex test proceed according to [complex-investigation-protocol](./references/complex-investigation-protocol.md).
6. Otherwise form a hypothesis on the cause of the issues.
7. Implement the fix.
8. Output the hypothesis and attempted fix, plus reasons why you think it would work.
9. Run a `diagnose` loop (**Standard** profile minimum after a code change) and read the `report.json` file to see if the fix works.
  Append to `diagnose-attempted-fixes-[test/package]-[flake/broken/timeout/slow].jsonl` file in this json format:
  ```json
  {"timestamp": "[current_timestamp]", "model": "[current-model] (e.g. `claude-sonnet-4.6/high`, `gemini-3.1-pro`)", "hypothesis": "Your original hypothesis for the issue", "experiment": "A concise summary of what you tried. Include small code snippets if helpful", "result": "Did it fix it or not? If not, give concise reason why", "next": "Next steps to attempt"}
  ```
10. GOTO 2
11. Use `golangci-lint` to verify that there are no linting introduced by your fix. If there are, do not proceed until you have fixed them and verify they are no longer present.

IF at any time the user interrupts or interjects during this loop, pick it up again where you left off, unless explicitly told otherwise.
</loop>

<tests-context>
Chainlink nodes are blockchain oracles. Read /README.md.
Tests share single postgres DB. Diagnose loop creates new DB.
</tests-context>

<flaky-test-flow>
Output hypothesis first. Show diffs. Do not abstract fixes.

Approaches:

1. Narrowing: Group failures. Ask user to proceed. Focus worst test otherwise.
2. Isolate: Pass alone, fail in package. Fix cross-test dependency.
3. Order: Shuffle alters pass rate. Fix cross-test leakage. Capture failing seed.
4. Race: Weird stack traces, nil pointers.
5. Timeout: Check logs for blocking ops, bad channel close, backpressure.
6. Resources: CI-only load failure. Check CPU, Mem. Use `go test` profiles (`-race`, `-cpuprofile`, `-trace`).
</flaky-test-flow>

<context_compaction>
Reference `diagnose-attempted-fixes-[test/package]-[flake/broken/timeout/slow].jsonl` when summarizing.
</context_compaction>

<possible_execution_issues>

- GOCACHE permissions sandbox error. STOP. Require user execution outside sandbox.
- Postgres `operation not permitted` sandbox error. STOP. Require user execution outside sandbox.
  </possible_execution_issues>

<logs_structure>
[resultsDir]/
|-- iteration-n.log.jsonl # Read only if needed. Full outputs.
|-- postgres-state-n.md # Read for DB error/hang. Final state.
|-- report.json # Read for summary. Extract args via `jq .run`.
|-- report.csv # DO NOT READ.
|-- logs/
|---- pkg_TestName_iter-n.log # Read for specific test failures.
</logs_structure>

<sub_agent_protocol>
1. Spawn `LogAnalyzer` when reading `logs/` or `iteration-n.log.jsonl`. Read ./references/log-analyzer-subagent.md.
2. Spawn `GithubFailureAnalyzer` when inspecting CI failure. Read ./references/github-failure-analyzer.md.
3. Spawn `JiraManager` when interacting with JIRA. Read ./references/jira-mananger-subagent.md.
</sub_agent_protocol>
