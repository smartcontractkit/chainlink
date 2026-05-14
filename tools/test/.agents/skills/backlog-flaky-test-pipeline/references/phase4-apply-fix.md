---
phase: phase4
model: sonnet
---

<phase id="phase4">

<purpose>
Apply and verify fixes for all PROCEED issues. Updates `ticket_records` with `applied_fix` and `verification` fields.
</purpose>

<substep id="4a" runs-in="parent">
<purpose>Ownership re-check before touching any files.</purpose>

For each PROCEED issue: call `mcp__atlassian__getJiraIssue`, verify assignee matches cached `accountId`.

If reassigned: report "KEY-NNN is now assigned to {displayName} — reach out before proceeding." Apply mid-flight abandonment rule (unassign + transition to Open + comment). Continue with remaining issues.
</substep>

<subagent id="verification" model="sonnet" instances="one per PROCEED issue that passed 4a" parallelism="single message">

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

4. **Verify compilation**:
   ```bash
   go build ./...
   ```
   If `go build` fails → fix the compilation error before proceeding. Do not revert — diagnose and correct. Only abandon if the fix approach is architecturally broken; explain why.

5. **Run the linter** (scoped to the changed package — never the whole repo):
   Derive `lint_scope` from `fix_file`: take the directory relative to the repo root and append `/...` (e.g. `fix_file = core/chains/evm/txmgr/foo_test.go` → `lint_scope = ./core/chains/evm/txmgr/...`).
   ```bash
   golangci-lint run {lint_scope}
   ```
   - **Lint passes** → set `lint_status = "ran"`, proceed.
   - **Lint finds violations that are fixable within the fix's scope** → fix them, re-run to confirm, set `lint_status = "ran"`, proceed.
   - **Lint finds violations that require changes outside the fix's scope** → set `lint_status = "failed"`, record violation summary in `lint_failure_detail`, return without fixing.
   - **Lint cannot execute** (binary missing, config error — not a lint violation) → set `lint_status = "skipped"`, record reason in `lint_failure_detail`, return without blocking.

6. **Rerun the test 10 times in independent processes**:
   ```bash
   # Go (detected via go.mod presence):
   for i in $(seq 10); do go test -race -shuffle=on -run "^{TestName}$" ./{package}/...; done
   # Adjust for non-Go projects based on detected language/build tool.
   ```

7. **Record result and return**:
   - **10/10 pass** → return `{ "verdict": "FIXED", "diff": "<git diff output>" }`.
   - **< 10/10 pass** → verdict `PARTIAL_FIX`:
     - Revert: `git restore {file}`.
     - Apply mid-flight abandonment rule: unassign, transition to "Open".
     - Write Investigation Update comment (OUTCOME = PARTIAL_FIX). "What was investigated": the suspected cause. "Hypothesis": `proposer_root_cause`. "What was tried": `fix_description` + attempted diff. "Why it didn't hold": test passed {n}/10 runs + first failure output (truncated to ~500 chars). "Recommended next step": `recommended_next_step` adapted as next direction, or N/A.
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
Write results into `ticket_records` (`applied_fix.diff` for FIXED, `verification.local_10x_passed`, `lint_status`, `lint_scope`).

For each ticket where `lint_status` is `"skipped"` or `"failed"`, gate on user decision before proceeding:

<user-prompt>
Lint {skipped | failed} for KEY-NNN.
{lint_failure_detail}

To run lint for just the changed package:
  golangci-lint run {lint_scope}

How would you like to proceed?
(a) Proceed anyway
(b) Wait — I'll fix lints myself (reply when ready to continue)
</user-prompt>

Record decision in `user_decisions`. In `--auto` mode: automatically choose (a) and log the lint status.
</on_subagent_return>

<on_complete>
Announce verdict for each issue: "Fix results: KEY-1 FIXED, KEY-2 PARTIAL_FIX (reverted and returned to queue)."

State: "Moving to commit and PR. Please review the fix files before confirming."

Read [phase5-commit-pr.md](phase5-commit-pr.md) and follow its instructions.
</on_complete>

</phase>
