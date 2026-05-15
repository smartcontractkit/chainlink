---
phase: phase4
model_tier: standard
---

<phase id="phase4">

<purpose>
Apply and verify fixes for all PROCEED issues. Updates `ticket_records` with `applied_fix` and `verification` fields.
</purpose>

<substep id="4a" runs-in="parent">
<purpose>Ownership re-check before touching any files.</purpose>

**Skip entirely in local mode** (`invocation.mode = "local"`) — no JIRA tickets are owned.

For each PROCEED issue (JIRA modes): follow `_shared-jira-flaky-ops/recheck-ownership.md` with `jira_key` and `accountId`. If result is `reassigned`: follow `_shared-jira-flaky-ops/abandon-ticket.md` for that ticket. Continue with remaining issues.
</substep>

<subagent id="verification" model_tier="standard" instances="one per PROCEED issue that passed 4a" parallelism="single message">

<inputs>
`key`, `fix_file`, `fix_line`, `fix_description`, `proposer_root_cause`, `recommended_next_step` (from phase3 result), `test_name`, `package` (from slim record), `accountId`, `cloudId`.
</inputs>

<steps>
1. **Library check**: read `go.mod` for libraries that already satisfy the fix's needs. Prefer them over writing bespoke implementations:
   - Retry logic → prefer `github.com/avast/retry-go` or another retry library already in `go.mod`.
   - Timeouts/context management → use `context` package + helpers already present.
   - Assertion helpers → use test helper libraries already imported in the test file.
   Only write new logic when no existing library covers the case.

2. **No tests for new helpers**: if the fix introduces a new helper function (retry wrapper, setup/teardown utility, wait helper), do **not** add a unit test for it. The only test that should run is the existing flaky test being fixed.

3. **Apply the code change.**

4. **Run the linter** (scoped to the changed package — never the whole repo; lint will catch compilation errors too):
   Derive `lint_scope` from `fix_file`: take the directory relative to the repo root and append `/...` (e.g. `fix_file = core/chains/evm/txmgr/foo_test.go` → `lint_scope = ./core/chains/evm/txmgr/...`).
   ```bash
   golangci-lint run {lint_scope}
   ```
   - **Lint passes** → set `lint_status = "ran"`, proceed.
   - **Lint finds violations that are fixable within the fix's scope** → fix them, re-run to confirm, set `lint_status = "ran"`, proceed.
   - **Lint finds violations that require changes outside the fix's scope** → set `lint_status = "failed"`, record violation summary in `lint_failure_detail`, return without fixing.
   - **Lint cannot execute** (binary missing, config error — not a lint violation) → set `lint_status = "skipped"`, record reason in `lint_failure_detail`, return without blocking.

6. **Run the chainlink `diagnose` tool to verify the fix** (10 iterations, parallelized):
   ```bash
   go -C tools/test run . diagnose --ai-output --iterations 10 --parallel-iterations 5 -- --run "^{TestName}$" --race --shuffle=on ./{package}/...
   ```
   Rules:
   - `--ai-output` is mandatory (machine-readable summary) and must appear **before** the `--` separator.
   - Harness flags (`--iterations`, `--parallel-iterations`, `--fail-fast-on`) go before `--`. `go test` flags go after `--`.
   - Reduce `--parallel-iterations` to `1` if the test holds external resources that don't tolerate concurrent runs (e.g. listens on a fixed port, opens a fixed temp file, or claims an exclusive lock). Otherwise default 5.
   - For additional flags: `go -C tools/test run . diagnose -h`.

   Parse the `--ai-output` summary to determine pass count (N out of 10). The verdict logic in step 7 is unchanged.

7. **Record result and return**:
   - **10/10 pass** → return `{ "verdict": "FIXED", "diff": "<git diff output>" }`.
   - **< 10/10 pass** → verdict `PARTIAL_FIX`:
     - Revert: `git restore {file}`.
     - **JIRA mode**: follow `_shared-jira-flaky-ops/abandon-ticket.md` (unassign + transition to Open), then follow `_shared-jira-flaky-ops/investigation-comment.md` to write `addCommentToJiraIssue` (OUTCOME = PARTIAL_FIX). "What was investigated": the suspected cause. "Hypothesis": `proposer_root_cause`. "What was tried": `fix_description` + attempted diff. "Why it didn't hold": test passed {n}/10 runs + first failure output (truncated to ~500 chars). "Recommended next step": `recommended_next_step` adapted as next direction, or N/A.
     - **Local mode**: no JIRA writes. Return `{ "verdict": "PARTIAL_FIX", "pass_count": N }` — include in final summary as `PARTIAL_FIX (reverted)`.
     - Return `{ "verdict": "PARTIAL_FIX", "pass_count": N }`.
</steps>

<output-schema>
```json
{
  "verdict": "FIXED | PARTIAL_FIX",
  "diff": "string (FIXED only — git diff of the applied change)",
  "pass_count": "integer (PARTIAL_FIX only)",
  "lint_status": "ran | skipped | failed",
  "lint_scope": "string — e.g. ./core/chains/evm/txmgr/...",
  "lint_failure_detail": "string | null — violation summary or execution error reason"
}
```
The parent never sees raw build output, lint output, or test logs — only the compact verdict.
</output-schema>

</subagent>

<on_subagent_return>
Write results into `ticket_records` (`applied_fix.diff` for FIXED, `verification.iterations_passed`, `lint_status`, `lint_scope`).

For each ticket where `lint_status` is `"skipped"` or `"failed"`, gate on user decision before proceeding:

Use `AskUserQuestion`:
- Question: "Lint {skipped | failed} for {key}. {lint_failure_detail}. To run lint manually: `golangci-lint run {lint_scope}`. How would you like to proceed?"
- Options: ["Proceed anyway", "Wait — I'll fix lints myself (reply when ready to continue)"]

Record decision in `user_decisions`. In `--auto` mode: automatically choose (a) and log the lint status.
</on_subagent_return>

<on_complete>
Announce verdict for each issue: "Fix results: KEY-1 FIXED, KEY-2 PARTIAL_FIX (reverted and returned to queue)." In local mode, use `local_id` or test name in place of the JIRA key.

- **JIRA mode**: State "Moving to commit and PR. Please review the fix files before confirming." → Read [phase5-commit-pr.md](phase5-commit-pr.md) and follow its instructions.
- **Local mode**: print the session summary:

```
Session complete (local mode — no JIRA, no PR).

Fix results:
| Test | Verdict | 10x | Notes |
|------|---------|-----|-------|
| <pkg>.TestFoo  | FIXED        | 10/10 | diff retained, uncommitted |
| <pkg>.TestBar  | PARTIAL_FIX  | 4/10  | reverted |
| <pkg>.TestBaz  | SKIPPED      | —     | classified SUT |
| <pkg>.TestQux  | INCONCLUSIVE | —     | debate did not converge |
| <pkg>.TestQuux | MISMATCH     | —     | stack trace stale |
```

Column rules:
- **Test**: `{package}.{test_name}` (use `local_id` as fallback if package is null).
- **Verdict**: FIXED | PARTIAL_FIX | SKIPPED | INCONCLUSIVE | MISMATCH.
- **10x**: pass count out of 10 for FIXED/PARTIAL_FIX; `—` for others.
- **Notes**: one short phrase. FIXED → "diff retained, uncommitted". PARTIAL_FIX → "reverted". SKIPPED → classification reason. INCONCLUSIVE → "debate did not converge". MISMATCH → "stack trace stale".

Then print: *"FIXED diffs are uncommitted in your working tree. Review with `git diff` and commit manually if you want to keep them."*
</on_complete>

</phase>
