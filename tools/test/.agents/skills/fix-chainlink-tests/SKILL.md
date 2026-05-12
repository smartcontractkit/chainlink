---
name: fix-chainlink-tests
description: >-
  Diagnoses and fixes unstable Chainlink Go tests (flakes, races, timeouts, deadlocks,
  slow runs). Use for non-deterministic failures or slow tests.
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
- DO NOT use plain `go test` commands. Only use `go -C tools/test run . diagnose`. Use `--iterations 1` for a single run.
- For `diagnose` runs expected >2m: Execute in background. Perform a single 30s crash check, then suspend task and wait for the report.json system notification. DO NOT poll.
</absolute_constraints>

<cli_reference>
Base command from `chainlink/` dir: `go -C tools/test run . diagnose [harness_flags] -- [go_test_flags] ./path`
From `chainlink/tools/test/` dir: `go run . diagnose [harness_flags] -- [go_test_flags] ./path` (same `./path` as above).
Package patterns after `--` are resolved against the monorepo `repo_root` (and may be rewritten when targeting a nested `go.mod`), not your shell cwd. Prefer `./core/...` / `./deployment/...`, not `../../...` from `tools/test/`.
- ALWAYS use `--ai-output` before the `--`.
- Harness flags (before `--`): `--iterations N`, `--fail-fast-on=(timeout|slow)`, `--parallel-iterations N`
- Go test flags (after `--`): `--run '^TestName$'`, `--timeout 10m`, `--race`
- Help: `go -C tools/test run . diagnose -h`
- Lint check: `golangci-lint run ./<packages-you-change> --fix`
</cli_reference>

<permissions-issues>
If you hit permissions issues running commands for this skill, read and/or direct the user to [permissions.md](/tools/test/.agents/skills/fix-chainlink-tests/references/permissions.md). If you are unable to quickly resolve these STOP. For the rest of the skill, ask the user to run permission-blocked commands for you.
</permissions-issues>

<init>
1. Verify target scope (test, package, or issue). If unknown, prompt user to provide one.
2. Check if you have access to the `trunk` MCP server. If not, prompt the user to [install and authenticate it](https://github.com/trunk-io/mcp-server#quick-start).
  <trunk>
  Trunk.io is a service we use to track flaky tests in CI. Access it with the Trunk MCP server. If available, use it to fetch more data on flaky or broken tests. **Use this as another data point in your diagnosis**, not as a definitive answer.
  1. Use `search-test` MCP tool to lookup find the test ID of the test(s) you're trying to fix.
  2. Use `fix-flaky-test` MCP tool to gather diagnosis data on a specific test.
  </trunk>
3. Investigate and understand the specific test code. 
4. Formulate initial hypothesis based on any data you can gather from the user or code.
</init>

<loop>
1. If a `diagnose-attempted-fixes-[test/package]-[flake/broken/timeout/slow].jsonl` file exists, read it to see previous fix attempts and findings. If not, create one.
2. Form a hypothesis on the cause of the issues
3. Implement a fix
4. Output the hypothesis and attempted fix, plus reasons why you think it would work.
5. Run a `diagnose` loop and read the `report.json` file using jq to see if the fix works. 
  Append to `diagnose-attempted-fixes-[test/package]-[flake/broken/timeout/slow].jsonl` file in this json format:
  ```json
  {"timestamp": "[current_timestamp]", "model": "[current-model] (e.g. `claude-sonnet-4.6/high`, `gemini-3.1-pro`)", "hypothesis": "Your original hypothesis for the issue", "experiment": "A concise summary of what you tried. Include small code snippets if helpful", "result": "Did it fix it or not? If not, give concise reason why", "next": "Next steps to attempt"}
  ```
6. If no issues, run a final `diagnose` loop to fully validate it's fixed. Use MINIMUM `100` iterations. Target `300-500` iterations for higher confidence.
7. If issues detected, focus on the ones the user wants to fix.

IF at any time the user interrupts or interjects during this loop, pick it up again where you left off, unless explicitly told otherwise.
</loop>

<tests-context>
* Chainlink nodes are blockchain oracles that read and write data to chains. Read the [README.md](/README.md) for more info.
* All tests share a single postgres DB. Each `diagnose` loop creates a new one.
</tests-context>

<analysis>
Lead with your hypothesis before writing code. Show contextual diffs, do not describe fixes abstractly. Common approaches and diagnoses:

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

- [filter.md](./references/flaky-patterns/filter.md): Tests using `Filter` functions to validate on-chain events. Usually LogPoller based tests. Commonly fails with: `failed to retrieve log value pointer after last block`
- [sql-lockout.md](./references/): `failed to create ...: ERROR: canceling statement due to lock timeout (SQLSTATE 55P03)`
</known_patterns>

<context_compaction>
When summarizing/compacting/compressing context:

- Strictly maintain a reference to the `attempted-fixes-[test/package]-[flake/broken/timeout/slow].jsonl` you're using for this session.
- Keep crucial data about your understanding of the test code you're trying to fix.
</context_compaction>

<logs_structure>
[resultsDir]/
|-- iteration-n.log.jsonl # DO NOT READ unless absolutely necessary; full log outputs, long and messy
|-- postgres-state-n.md # Final state of tests' postgres DB after iteration. Read if diagnosing DB-based errors or hangs.
|-- report.json # Read this; summary of full `diagnose` run (include `jq .run` for go test args and harness flags)
|-- report.csv # DO NOT READ; human readable csv
|-- logs/ # Extracted individual test logs
|---- pkg_TestName_iter-n.log # Logs for individual slow/failing tests, read this as needed
</logs_structure>

<read_logs_sub_agent>
When reading log files from the `logs/` directory or `iteration-n.log.jsonl`, you MUST spawn a specialist `LogAnalyzer` sub-agent.

You MUST configure the sub-agent with these exact initialization parameters:

1. System Prompt: "You are a headless, read-only log parser. Your sole purpose is to read Go test logs from the end up. Each log file contains logs from `chainlink` nodes, plus test-specific logs. Read the logs and construct possible reasons why the test [input reason we're investigating]. You do not converse. You output raw JSON and nothing else."
2. Allowed Tools: File read/grep tools ONLY. Revoke all execution, write, and web search capabilities.
3. Temperature: 0.0

The sub-agent MUST output ONLY valid JSON matching this exact structure. DO NOT wrap the output in markdown code blocks. Output raw JSON only, with no explanations and no yapping:

```json
{
  "logs_read": ["log_path_1.log", "log_path_2.log"],
  "failure_diagnosis": [
    {
      "possible_reason": "reason for failure",
      "evidence": "specific logs/log lines"
    },
    {
      "possible_reason": "reason for failure",
      "evidence": "specific logs/log lines"
    }
  ]
}
```
</read_logs_sub_agent>
