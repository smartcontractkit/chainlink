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

<context_compaction>
When summarizing context, strictly maintain state in this format:

## [TestName]
Failure: [suspected failure reasons]
SuspectedFix: [the fix you've implemented or want to try]
NextStep: [the next step for diagnosing/fixing/verifying the test]
</context_compaction>

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
- Shuffle test order: `go test -shuffle=on -count=50 -failfast ./path/to/package`
- CPU/Memory load: `go test -cpu=1,2,4 -count=20 -failfast ./path/to/package`
- Lint check: `golangci-lint run ./<packages-you-change> --fix`
</cli_reference>

## Execution & Analysis
- **Postgres:** Serial diagnose restores DB between iterations. Parallel gives each worker an ephemeral DB. Neither resets between tests *within* one iteration.
- **Report Analysis:** Read `<resultsDir>/report.json` using `jq`. Top-level buckets: `flakes`, `failures`, `timeouts`, `slow`. Harness and `go test` invocation: `jq .run` (argv, iteration count, fail-fast, shuffle, etc.).
- **Narrowing:** If many tests flag, look for similarities in their failures. If found, present that to the user and ask if they want to continue with that assumption. If not, try to focus on the most problematic test.
- **Profiles:** When logs/report are insufficient, use standard `go test` profile flags (`-race`, `-cpuprofile`, `-trace`, etc.). View with `go tool pprof` or `go tool trace`.

<logs_structure>
<resultsDir>/
|-- iteration-n.log.jsonl # DO NOT READ unless absolutely necessary; full log outputs, long and messy
|-- postgres-state-n.md # Final state of postgres DB after test iteration. Read if diagnosing DB-based errors or hangs.
|-- report.json # Read this; summary of full `diagnose` run (include `jq .run` for go test args and harness flags)
|-- report.csv # DO NOT READ; human readable csv
|-- logs/ # Extracted individual test logs
|---- pkg_TestName_iter-n.log # Logs for individual slow/failing test
</logs_structure>

<sub_agent_protocol>
When reading log files from the `logs/` directory or `iteration-n.log.jsonl`, you MUST spawn a sub-agent to read from the end up. 
The sub-agent MUST output ONLY valid JSON matching this exact structure, with no markdown, no explanations, and no yapping:
{
  "logs_read": ["log_path_1.log", "log_path_2.log"],
  "failure_diagnosis": [
    {
      "possible_reason": "explanation",
      "evidence": "reasoning and evidence"
    }
  ]
}
</sub_agent_protocol>

## Playbook & General Fixes
Lead with your hypothesis before writing code. Show contextual diffs, do not describe fixes abstractly.

1. **Check Known Patterns:** See `<known_patterns>` below for common flaky test patterns and fixes in this repo. Try them first.
2. **Isolate (Pass alone, fail in package):** Cross-test dependency. Missing `t.Cleanup`, global state (`var` singletons, loggers), or shared mock servers. Fix by moving state to per-test constructors or using `t.Cleanup`.
3. **Order (Shuffle changes pass rate):** Same as isolation. Fix cross-test leakage. Capture failing seed and provide to user.
4. **Race:** Triggers on weird stack traces or nil pointers. Use `-race`. Fix with `sync.Mutex`, `atomic.*`, or narrow shared fields.
5. **Timeout:** Check logs for blocking (chan receive, `Wait`, `testutils.WaitTimeout`). Use `synctest` to improve tests relying on channels.
6. **Slow:** Compare `p50` vs `max_elapsed`. Look for `time.Sleep` or coarse polling loops. Replace with `require.eventually` or channel sync. Simulated chains are frequent offenders.
7. **Resources:** If failing under load/CI only, DB connections might be exhausted by `t.Parallel()`. Use separate schema/user per test. 

<known_patterns>
  <pattern name="LogPoller Timing Race">
    <symptom>
      The dominant flake pattern in simulated-chain tests that enable `Feature.LogPoller = true`. Error message contains `"failed to retrieve log value pointer of block N: not found"` and the stack trace points to a `FilterXxx` call that immediately follows a `backend.Commit()`. Note: Raw geth bindings do NOT have this race, only interface types backed by LogPoller.
    </symptom>
    
    <fix_a_receipt_parsing>
      For one-shot events where you only need a value emitted at creation (e.g. `SubscriptionCreated`, `RequestSent`): parse the tx receipt directly instead of calling `FilterXxx`.
      ```go
      // AFTER (deterministic):
      tx, err := coordinator.CreateSubscription(auth)
      require.NoError(t, err)
      backend.Commit()
      receipt, err := backend.Client().TransactionReceipt(ctx, tx.Hash())
      require.NoError(t, err)
      require.Equal(t, uint64(1), receipt.Status)
      var subID *big.Int
      for _, log := range receipt.Logs {
          if log.Address != coordinatorAddress {
              continue
          }
          // SubscriptionCreated(uint64 indexed subId, address owner): Topics[1] = subId
          subID = new(big.Int).SetBytes(log.Topics[1].Bytes())
          break
      }
      require.NotNil(t, subID, "no SubscriptionCreated log in receipt")
      ```
    </fix_a_receipt_parsing>

    <fix_b_non_fatal_filter>
      For diagnostic/verification filters called inside a polling loop: a transient LogPoller error must not crash the test — it should retry.
      ```go
      // AFTER (retries):
      require.Eventually(t, func() bool {
          // LogPoller may not have indexed the latest block yet; skip and retry.
          it, err := coordinator.FilterRandomWordsForced(nil, ids, subs, addrs)
          if err == nil {
              for it.Next() {
                  require.Equal(t, expected, it.Event.Field)
              }
          }
          return utils.IsEmpty(commitment[:])
      }, timeout, tick)
      ```
    </fix_b_non_fatal_filter>

    <fix_c_dynamic_reference>
      If `require.Eventually` commits new blocks on each iteration, compute the reference block number inside the closure so it doesn't become stale.
      ```go
      // AFTER (dynamic):
      require.Eventually(t, func() bool {
          backend.Commit()
          tip, err := backend.Client().HeaderByNumber(ctx, nil)
          if err != nil || tip == nil || tip.Number.Uint64() < 256 {
              return false
          }
          _, err = bhsContract.GetBlockhash(nil, new(big.Int).SetUint64(tip.Number.Uint64()-256))
          return err == nil
      }, testutils.WaitTimeoutCustom(t, 5*time.Minute), time.Second)
      ```
    </fix_c_dynamic_reference>
  </pattern>

  <pattern name="TXM broadcast latency (parallel load)">
    <symptom>
      Under 5+ parallel test workers, TXM broadcasts transactions asynchronously. A heartbeat/fulfillment tx may be logged as "sent" by the service but not yet in the mempool when the next `backend.Commit()` fires. Test detects service as active, but stored block is `N+1` or later than the fixed reference.
    </symptom>
    <fix>
      Use the dynamic reference fix (`fix_c_dynamic_reference` from LogPoller Timing Race) so the check tracks wherever the tx actually lands.
    </fix>
  </pattern>
</known_patterns>