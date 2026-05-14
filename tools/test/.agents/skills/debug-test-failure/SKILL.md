---
name: debug-test-failure
description: >-
  A deep-dive diagnostic tool for fixing Go test failures (flakes, races, timeouts,
  deadlocks) identified during local development or active CI failures.

  USE THIS WHEN:
  1. You have a specific, known failing test name or local error log.
  2. You are currently working on a branch and need to fix a regression or a new flake.
  3. The failure is NOT necessarily tracked in JIRA or indexed by Trunk.io yet.
  4. You need to perform deep "forensic" code analysis and manual fix iterations.

  DO NOT USE THIS WHEN:
  1. Processing a batch of existing JIRA tickets (use 'backlog-flaky-test-pipeline' instead).
  2. You require automated JIRA status updates or Trunk.io integration.
  3. The goal is bulk technical debt reduction from a project backlog.
---

<absolute_constraints>
- DO NOT use this skill if the user already has a known fix (apply it directly).
- DO NOT use for deterministic first-run failures (use normal debug).
- DO NOT use for full-suite CI prep (use `make new_test` or `make new_gotestsum` instead).
- ONLY run tests in these packages without explicit user approval: `core/`, `deployment/`. Warn the user if running outside these.
- DO NOT modify the test's core goal to make it pass.
- DO NOT remove tests/assertions unless replacing with better ones or deleting confirmed dead code.
- DO NOT modify package-wide helpers (`testutils`) to fix localized tests.
- DO NOT use plain `go test` commands. Only use `go -C tools/test run . diagnose`. Use `--iterations 1` for a single run.
- For `diagnose` runs expected >2m: Execute in background. Perform a single 30s crash check, then suspend task and wait for the report.json system notification. DO NOT poll.
</absolute_constraints>

## Initialization
1. Verify target scope (test, package, or issue). If unknown, prompt user.
2. Formulate initial hypothesis: flake, timeout, slow, panic, deadlock, or race.
3. Run bounded diagnosis (`--fail-fast` or low `--iterations`).

<cli_reference>
Base Command: `go -C tools/test run . diagnose [harness_flags] -- [go_test_flags] ./path`
- ALWAYS use `--ai-output` before the `--`.
- Harness flags (before `--`): `--iterations N`, `--fail-fast-on=(timeout|slow)`, `--parallel-iterations N`
- Go test flags (after `--`): `--run '^TestName$'`, `--timeout 10m`, `--race`
- Help: `go -C tools/test run . diagnose -h`
- Lint check: `golangci-lint run ./<packages-you-change> --fix`
</cli_reference>

<loop>
1. If user doesn't have recent results, run `diagnose` command with min 5 iterations to gather initial info. If issues due to sandbox, STOP, give user command to run and have them report results.
2. If no issues, ask the user if they want to verify with more iterations. If not, end and output final report of findings, fixes, and lessons learned.
3. If issues detected, focus on the ones the user wants to fix.
4. If a `diagnose-attempted-fixes-[test/package]-[flake/broken/timeout/slow].jsonl` file exists, read it to see previous fix attempts and findings.
5. Form a hypothesis on the cause of the issues
6. Implement a fix
7. Output the hypothesis and attempted fix, plus reasons why you think it would work.
8. Run a `diagnose` loop and read the `report.json` file using jq to see if the fix works.
  Append to `diagnose-attempted-fixes-[test/package]-[flake/broken/timeout/slow].jsonl` file in this json format:
  ```json
  {"timestamp": "[current_timestamp]", "model": "[current-model] (e.g. `claude-sonnet-4.6/high`, `gemini-3.1-pro`)", "hypothesis": "Your original hypothesis for the issue", "experiment": "A concise summary of what you tried. Include small code snippets if helpful", "result": "Did it fix it or not? If not, give concise reason why", "next": "Next steps to attempt"}
  ```
9. GOTO 2

IF at any time the user interrupts or interjects during this loop, pick it up again where you left off, unless explicitly told otherwise.
</loop>

<tests-context>
* Chainlink nodes are blockchain oracles. Read the [README.md](/README.md)
* All tests share a single postgres DB. Each `diagnose` loop creates a new one.
</tests-context>

<analysis>
Lead with your hypothesis before writing code. Show contextual diffs, do not describe fixes abstractly. List of common approaches and diagnoses:

1. **Check Known Patterns:** See `<known_patterns>` below for common flaky test patterns and fixes in this repo. If they apply to the situation attempt them first.
2. **Narrowing:** If many tests flag, look for similarities in their failures. If found, present that to the user and ask if they want to continue with assumption of relation. If not, try to focus on the most problematic test.
3. **Isolate (Pass alone, fail in package):** Cross-test dependency. Missing `t.Cleanup`, global state (`var` singletons, loggers), or shared mock servers. Fix by moving state to per-test constructors or using `t.Cleanup`.
4. **Order (Shuffle changes pass rate):** Same as isolation. Fix cross-test leakage. Capture failing seed and provide to user.
5. **Race:** Triggers on weird stack traces or nil pointers. Use `-race`. Fix with `sync.Mutex`, `atomic.*`, or narrow shared fields.
6. **Timeout:** Check logs for blocking (chan receive, `Wait`, `testutils.WaitTimeout`). Use `synctest` to improve tests relying on channels.
7. **Slow:** Compare `p50` vs `max_elapsed`. Look for `time.Sleep` or coarse polling loops. Replace with `require.Eventually` or channel sync. Simulated chains are frequent offenders.
8. **Resources:** If failing under load/CI only, check CPU and Memory usage. When logs/report are insufficient, use standard `go test` profile flags (`-race`, `-cpuprofile`, `-trace`, etc.). View with `go tool pprof` or `go tool trace`.
</analysis>

<known_patterns>
Files in the `references/flaky-patterns/` dir.
- [filter.md](./references/flaky-patterns/filter.md): Tests using `Filter` functions to validate on-chain events. Usually LogPoller based tests.
- [sql-lockout.md](./references/): `failed to create ...: ERROR: canceling statement due to lock timeout (SQLSTATE 55P03)`
</known_patterns>

<context_compaction>
When summarizing/compacting/compressing context, strictly maintain a reference to the `attempted-fixes-[test/package]-[flake/broken/timeout/slow].jsonl` you're using for this session.
</context_compaction>

<possible_execution_issues>
- **GOCACHE permissions issues**: `[build failed]\n open .../Library/Caches/...` This is caused by some sandbox environments. If you cannot exit the sandbox to fix this, STOP. DO NOT attempt to create a new cache. Ask the user to run the command instead and give you results so you can continue.
- **Postgres sandbox error**: `operation not permitted` connecting to postgres. Sandbox issues. If you cannot exit the sandbox to fix this, STOP. Ask the user to run the command instead and give you results so you can continue.
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
When reading log files from the `logs/` directory or `iteration-n.log.jsonl`, you MUST spawn a specialist `LogAnalyzer` sub-agent.

You MUST configure the sub-agent with these exact initialization parameters:
1. System Prompt: "You are a headless, read-only log parser. Your sole purpose is to read Go test logs from the end up. Each log file contains logs from `chainlink` nodes, plus test-specific logs. Read the logs and construct possible reasons why the test [input reason we're investigating]. You do not converse. You output raw JSON and nothing else."
2. Allowed Tools: File read/grep tools ONLY. Revoke all execution, write, and web search capabilities.
3. Temperature: 0.0

The sub-agent MUST output ONLY valid JSON matching this exact structure. DO NOT wrap the output in markdown code blocks. Output raw JSON only, with no explanations and no yapping:
{
  "logs_read": ["log_path_1.log", "log_path_2.log"],
  "failure_diagnosis": [
    {
      "possible_reason": "explanation",
      "evidence": "specific logs/log lines"
    }
  ]
}
</sub_agent_protocol>
