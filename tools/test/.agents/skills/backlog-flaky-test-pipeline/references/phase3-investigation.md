---
phase: phase3
model: sonnet
---

<phase id="phase3">

<purpose>
Run parallel per-ticket investigations: extract Trunk data, analyze test code, classify flakiness source, debate root cause, and return a structured verdict per ticket. Updates `ticket_records` with `actionable_facts`, `chosen_fix`, and outcome. Runs a checkpoint gate before any files are modified.
</purpose>

<prereqs>

Read [shared-jira-protocol.md](shared-jira-protocol.md) if not already loaded in this session. Required for Investigation Update comment format and mid-flight abandonment procedure.
</prereqs>

<parallelism>
Spawn all N per-ticket investigation subagents in a **single message**. Each receives: the slim record (`key`, `title`, `description`, `trunk_test_case_url`, `test_name`, `package`, `previous_attempts`), `nav_tool`, `lsp_available`, `auto_mode`. Never pass raw JIRA API objects.
</parallelism>

---

<substep id="3a" model="haiku">
<purpose>Resolve `testCaseId`, retrieve `fix-flaky-test` historical data, optionally invoke `investigate-ci-failure` for single-run forensic data, and run the top-level subtest check.</purpose>

<fix-flaky-test-call>
- Use `slim_record.test_case_id` as `testCaseId`. If null: call `mcp__trunk__search-test` with `slim_record.test_name` (falling back to `slim_record.title`) as fuzzy fallback — flag if used, it may match the wrong test.
- Call `mcp__trunk__fix-flaky-test` with `testCaseId` (no `createNewInvestigation`). **Immediately after each response** (including polls), apply the filter — never store or forward the raw blob:
  1. Capture `trunk_analysis_url` from `investigationUrl` if present.
  2. Extract only facts with `Confidence ≥ 0.9` from `## Facts`. Discard all lower-confidence facts and the entire `## Markdown Summary`.
  3. Store as `trunk_filtered_facts: [str]`.

**If response has `summary` or `facts`**: set `trunk_investigation_status = "existing"`.

**If response is empty** (no investigation yet) — branch on whether a CI URL is available:

  **Branch A: `slim_record.ci_run_url` is non-null** → skip polling, rely on `investigate-ci-failure` (next block) instead. Set `trunk_investigation_status = "ci_run_only"`.

  **Branch B: `ci_run_url` is null AND `auto_mode` is false** → opt-in fallback prompt:

  <user-prompt>
  KEY-NNN: No existing Trunk investigation found. If you have a recent CI failure URL for this test, paste it now to analyze that run directly. Otherwise press enter to trigger a new investigation (this takes 2–5 minutes).
  </user-prompt>

  - URL provided → set `slim_record.ci_run_url` and proceed to `<investigate-ci-failure-call>`. Set `trunk_investigation_status = "ci_run_only"`.
  - User pressed enter (no URL) → fall through to **Branch C**.

  **Branch C: no CI URL and (auto_mode or user declined fallback)** → existing flow:
  - Inform user: *"No existing Trunk investigation found for {KEY}. Triggering one now — this may take 2–5 minutes."*
  - Call `mcp__trunk__fix-flaky-test` with `createNewInvestigation: true`.
  - Poll every ~30 seconds for up to 5 minutes (10 attempts). Apply the same immediate filter on each response.
  - If `investigationId` is unknown during polling: expected — the investigation is still initializing. Inform user: *"Investigation pending (investigationId not yet active — this is normal). Continuing to poll…"*
  - `trunk_filtered_facts` becomes non-empty → set `trunk_investigation_status = "triggered"`.
  - Still empty after 5 minutes → set `trunk_investigation_status = "uninvestigated"`. Ask user:

  <user-prompt>
  Trunk did not return investigation results after 5 minutes for {KEY}. How would you like to proceed?
  (a) Continue with code analysis only  (b) Wait longer  (c) Skip this issue
  </user-prompt>
</fix-flaky-test-call>

<investigate-ci-failure-call>
Runs only when `slim_record.ci_run_url` is non-null (after upfront input or fallback prompt).

- Call `mcp__trunk__investigate-ci-failure` with `workflowUrl = slim_record.ci_run_url` and `orgSlug = "chainlink"` (hardcoded — see `trunk-org-slug` tip).
- Store the full structured response as `ci_run_evidence` on the ticket record. **No `Confidence ≥ 0.9` filter applies** — this is a separate evidence track per the constraints exemption.
- Edge cases:
  - Tool returns "build/compile failure before tests ran" → set `ci_run_evidence = { "status": "build_failure", "raw": <response> }`. Inform user: *"CI run for {KEY} failed before tests executed; CI evidence will be limited."* Continue.
  - Tool errors (permission, unknown error, malformed URL) → set `ci_run_evidence = null`, log the error to the user, do not retry.
- The 3b-ii classifier does **not** consume `ci_run_evidence` — only the Proposer does (see 3d).
</investigate-ci-failure-call>

<top-level-check>
After Trunk investigation resolves, inspect `trunk_filtered_facts` (and any stack trace within them) for a `file:line`. If the failure line falls inside a `t.Run(...)` callback AND the outer function contains no assertions outside `t.Run` blocks → candidate for **SKIP_TOP_LEVEL**.

**Exception**: if `slim_record.test_name` contains a `/` the ticket was already filed against a specific subtest — SKIP_TOP_LEVEL must not fire. If `test_name` is null, check the title for `/`.

If all conditions are met:
1. `mcp__atlassian__getTransitionsForJiraIssue` → find a closing transition ("Won't Do", "Closed", "Done"). If it supports a resolution field, set `resolution = "Won't Do"` (fallback: "Won't Fix").
2. `mcp__atlassian__transitionJiraIssue` → close with that transition + resolution.
3. `mcp__atlassian__addCommentToJiraIssue` → Investigation Update comment (OUTCOME = CLOSED_SUBTEST). "What was investigated": failure originates in a `t.Run` subtest, not the top-level function. "Recommended next step": file or locate a ticket for the specific subtest. All other sections: N/A.
4. Stop investigation for this issue.
</top-level-check>

</substep>

---

<substep id="3b" parallelism="Subagent A and Subagent B spawn in a single message">

<subagent id="trunk-analyzer" model="haiku">
<inputs>
`trunk_filtered_facts` (already filtered to ≥ 0.9 in 3a), `trunk_investigation_status`. Do not call any Trunk MCP tools.
</inputs>
<logic>
- If `trunk_investigation_status = "uninvestigated"` or `trunk_filtered_facts` is empty → return `confidence: "none"` with empty `facts`.
- Map confidence from pre-filtered facts:
  - `"high"`: at least one fact contains raw CI observational data — exact log lines, error messages, stack traces, or specific `file:line` from actual failing runs.
  - `"low"`: facts describe symptoms only (e.g. "failures in cluster 2 are regex mismatches") but contain no raw CI data.
  - `"none"`: `trunk_filtered_facts` empty or `trunk_investigation_status = "uninvestigated"`.
</logic>
<output-schema>
```json
{ "testCaseId": "string", "facts": ["string"], "confidence": "high | low | none" }
```
`facts` must be an array of raw text strings. Empty array is valid only when `confidence = "none"`.
</output-schema>
</subagent>

<subagent id="code-analyzer" model="sonnet">
<inputs>
Slim record, `nav_tool`, `lsp_available`, `trunk_filtered_facts`.
</inputs>
<logic>
**Locate test file** — never use grep/find if a smarter tool is available. Use the top-level function from `slim_record.test_name` (part before first `/`); scope with `slim_record.package` if available. Fall back to extracting from `slim_record.title` only if `test_name` is null.
- `nav_tool = "lsp"` or `lsp_available`: LSP definition lookup → extract `uri` + `range.start`.
- `nav_tool = "crg"`: `mcp__code-review-graph__semantic_search_nodes_tool` or `query_graph_tool` with `callees_of`.
- Last resort: `grep -r "func {TestName}" .`; parse first `filepath:line`; warn if multiple matches.

**Stacktrace currency check** — if `trunk_filtered_facts` contains a stack trace, find the innermost frame that belongs to test code (deepest frame whose file path is inside the repo, not vendor/framework). Check whether that function still exists using LSP preferred, then code-review-graph or grep as fallback.
- Function completely absent → `code_mismatch: true`, record in `mismatch_details`.
- Function exists (even at different line) → `code_mismatch: false`.
- No stack trace present → `code_mismatch: false`.

**Code analysis**: read the test and its helpers. Analyze for: timing dependencies, shared global state, ordering assumptions, missing cleanup/teardown, non-deterministic data.

**Previous attempts constraint**: if `previous_attempts` is non-empty, list each attempt explicitly and exclude those approaches. State: *"Previously attempted (will not re-propose): [X, Y]. Rejected because: [rejection_reasons]."* Surface any non-null `recommended_next_step` from prior attempts.

**Parallelism bias guard**: do not name parallel execution as `suspected_cause` unless you can identify a specific shared resource (global variable, shared file, database table, network port) written by one execution and read by another without synchronization. "Tests may run in parallel" is not a valid suspected cause.
</logic>
<output-schema>
```json
{
  "file": "string",
  "line": "integer",
  "analysis": "string",
  "suspected_cause": "string | null",
  "suspected_cause_location": "test_code | production_code | unknown",
  "excluded_approaches": ["string"],
  "code_mismatch": "boolean",
  "mismatch_details": "string | null"
}
```
</output-schema>
</subagent>

<schema-validation applies-to="trunk-analyzer code-analyzer">
Two failure classes — validate each subagent individually:

- **Transient** (empty or null response — likely tool timeout): retry immediately with original prompt.
- **Structural** (fields missing, wrong types, semantically invalid): retry with (1) exact validation error, (2) concrete example of expected format for the failing field, (3) subagent's previous invalid output.

Structural failures include: `facts` containing category labels or counts instead of raw text strings (`"CI_LOGS (1.0)"` is a label — invalid; `"Error: no contract code at given address"` is raw text — valid), `confidence` not one of three allowed values, `code_mismatch` not a boolean.

Allow up to **3 total attempts** per subagent. After 3 failures → hard stop for this issue. Apply mid-flight abandonment rule. Write Investigation Update comment (OUTCOME = ABANDONED): state which subagent failed and include the validation error; recommended next step: re-run and include last raw output verbatim. Continue with other issues.
</schema-validation>

</substep>

---

<substep id="3b-ii" model="haiku→sonnet" runs-in="parent">
<purpose>
Classify flakiness source as TEST / SUT / INFRA / AMBIGUOUS before entering the fix debate. Runs in the parent (not a subagent) because it may require a user gate.
</purpose>

<input-schema>
```json
{
  "$schema": "phase_3bii_input_v1",
  "ticket_key": "string",
  "test_name": "string",
  "trunk_filtered_facts": ["string"],
  "trunk_investigation_status": "existing | triggered | uninvestigated | ci_run_only",
  "subagent_b_output": {
    "file": "string", "line": "number", "analysis": "string",
    "suspected_cause": "string", "suspected_cause_location": "test_code | production_code | unknown",
    "excluded_approaches": ["string"], "code_mismatch": "boolean", "mismatch_details": "string | null"
  }
}
```
**Matching scope rule**: all signal triggers match only against `trunk_filtered_facts` text or stack-trace excerpts — never against test source code, test names, or code comments.
</input-schema>

<tier id="1" type="string-match">
Deterministic — no LLM judgment.

**SUT signals:**
| ID | Trigger | Weight |
|----|---------|--------|
| SUT_SERVICE_UNAVAILABLE | "connection refused", "service unavailable", "dial tcp", "reset by peer"; OR "EOF" co-occurring within 200 chars of any of {grpc, rpc, dial, net., tcp, http} | 2 |
| SUT_COMPONENT_NOT_INITIALIZED | "nil pointer dereference", "not initialized", "component not ready" — excerpt must be in a production code frame (not `_test.go`) | 1 |
| SUT_CONSISTENT_PROD_CODE_FAILURE | Same production-code `file:line` appears in stack traces from ≥ 2 distinct CI runs in `trunk_filtered_facts` | 1 |

**TEST signals** (all weight 1):
| ID | Trigger |
|----|---------|
| TEST_SHARED_GLOBAL_STATE | Package-level var, global map, or singleton mutated without `t.Cleanup` restore |
| TEST_PARALLEL_UNSYNC | `t.Parallel()` present with shared resource used without sync primitive |
| TEST_TIMING_DEPENDENCY | `time.Sleep` or fixed-duration delay used to synchronize async behavior |
| TEST_MISSING_CLEANUP | Resource setup (server, DB, goroutine) without corresponding Cleanup/defer |
| TEST_RACE_DETECTOR_FIRED | `DATA RACE` or `RACE CONDITION DETECTED` in `trunk_filtered_facts` |

**INFRA signals** (any match → INFRA, overrides scoring):
| ID | Trigger |
|----|---------|
| INFRA_OOM_KILLED | "signal: killed", "OOM", "out of memory", "exit status 137" |
| INFRA_DISK_FULL | "no space left on device", "disk full" |
| INFRA_REGISTRY_FAILURE | "pulling image", "registry", "manifest unknown", "pull access denied" |
</tier>

<tier id="2" type="semantic" model="sonnet">
| ID | Trigger | Weight |
|----|---------|--------|
| SUT_PRECONDITION_NOT_MET | LLM determines the failure indicates a precondition was unsatisfied — a dependency wasn't available, a registration hadn't completed, a service hadn't started — rather than the SUT behaving incorrectly *during* the test scenario. LLM must quote a verbatim excerpt from `trunk_filtered_facts` and provide a one-sentence explanation. | 2 |

Ask: *"Does this failure indicate that a precondition for the test wasn't met (something wasn't ready or available), or does it indicate the system-under-test behaved incorrectly while executing the test's actual scenario?"* If LLM answers yes with a verbatim excerpt, the signal fires. If no excerpt can be quoted from the inputs, the signal is dropped.
</tier>

<scoring>
Deterministic post-LLM — LLM never computes scores:
1. LLM outputs matched signal IDs with one verbatim excerpt each.
2. Post-processing: each excerpt must appear verbatim in `trunk_filtered_facts` joined text OR `subagent_b_output.analysis`. Unvalidated signals are dropped in-place; scores recomputed (no retry).
3. `sut_score` = sum of weights for validated SUT signals; `test_score` = sum of weights for validated TEST signals.
4. Classify: any validated INFRA signal → `INFRA`; `sut_score > test_score` → `SUT`; `test_score > sut_score` → `TEST`; equal → tiebreaker.

**Tiebreaker** (only on tie):
0. `sut_score == 0 AND test_score == 0` → `AMBIGUOUS` immediately.
1. `subagent_b_output.code_mismatch == true` → `AMBIGUOUS` (stale data, cannot trust evidence).
2. `trunk_investigation_status == "uninvestigated"` → `AMBIGUOUS`.
3. Any SUT signal excerpt verbatim-matches a SUT trigger string → `SUT`.
4. Default → `AMBIGUOUS`.

**Confidence rule**: `high` = score margin ≥ 2 OR at least one weight-2 signal on the winning side; `low` = margin = 1 with no weight-2 signal; `none` = no signals matched, or classification is AMBIGUOUS/INFRA.
</scoring>

<output-schema>
```json
{
  "$schema": "phase_3bii_output_v1",
  "classification": "TEST | SUT | AMBIGUOUS | INFRA",
  "confidence": "high | low | none",
  "sut_score": "number", "test_score": "number",
  "sut_signals_matched": ["string"], "test_signals_matched": ["string"], "infra_signals_matched": ["string"],
  "evidence": [{ "signal_id": "string", "source": "trunk_fact | code_analysis | stack_trace", "excerpt": "string" }],
  "tiebreaker_applied": "boolean", "tiebreaker_step_fired": "number | null",
  "rationale": "string",
  "sut_description": "string | null",
  "sut_pivot": { "file": "string | null", "component": "string | null", "hypothesis": "string | null" },
  "smell_notes": ["string"]
}
```
`sut_pivot`: required (fields may be null) when classification is SUT or AMBIGUOUS; null otherwise.
</output-schema>

<schema-validation applies-to="3b-ii">
- **Transient**: retry immediately with original prompt.
- **Structural** (missing fields, wrong types, invalid enum, `sut_pivot` absent when classification is SUT/AMBIGUOUS): retry with error context.
- Excerpt validation is post-processing, not a retry trigger — drop unvalidated signals and recompute.
- Allow up to 3 total attempts. After 3 failures: set `classification = "AMBIGUOUS"`, `confidence = "none"`, `rationale = "Schema validation failed after 3 attempts"`, and continue to gate logic.
</schema-validation>

<gate-logic>
| Classification | `--auto` mode | Interactive mode |
|---|---|---|
| TEST | Continue to 3c | Continue to 3c |
| SUT | Return to queue + JIRA comment | Prompt user (options a/b/c) |
| AMBIGUOUS | Return to queue + JIRA comment | Prompt user with both signal lists |
| INFRA | Return to queue + JIRA comment | Prompt user |

<user-prompt id="sut-gate">
Classification: **SUT** (confidence: {confidence})
Signals: {sut_signals_matched}
{sut_description}

This test appears to expose a SUT bug, not a test-code bug. Options:
(a) Return this ticket to the queue with a SUT-hypothesis note
(b) Override — treat as TEST and proceed to fix debate (audited)
(c) Fix the test code AND auto-file a SUT bug ticket (label: sut-bug)
</user-prompt>

*Option (b) audit trail*: add JIRA comment "Classification overridden to TEST by user. Original: SUT, confidence: {confidence}, signals: {list}." Add commit trailer: `Flakiness-classification: TEST (user override from SUT)`.

*Option (c)*: create a JIRA issue in the same project: summary `SUT bug: {sut_description}`, description includes `sut_pivot` fields, label `sut-bug`. Return the new ticket key to the user. Proceed to 3c treating current ticket as TEST.

For SUT/AMBIGUOUS/INFRA auto-queue returns: write Investigation Update comment (OUTCOME = RETURNED_TO_QUEUE). "What was investigated": classification, confidence, matched signal IDs, SUT score, TEST score, rationale. "Hypothesis": `sut_description` if SUT, otherwise N/A. "Recommended next step": SUT → investigate `sut_pivot`; AMBIGUOUS → clarify classification before re-investigating; INFRA → check infrastructure. Apply mid-flight abandonment rule (unassign + transition to Open). Continue with other issues.
</gate-logic>

</substep>

---

<substep id="3c" model="sonnet">
<purpose>Mismatch short-circuit — check before entering the debate.</purpose>

If `code_mismatch: true` from Subagent B → stop here. Return verdict `MISMATCH` with `mismatch_details`. Do not enter the debate — the Trunk stack trace references code that no longer exists; any fix derived from it would target the wrong location.

Otherwise: proceed to 3d. Trunk facts are seed evidence for the Proposer, not a replacement for code analysis and debate. High-confidence Trunk data tells you *what* failed in CI — it cannot tell you *why* the code produces that output or *how* to fix it. Code analysis is always required.
</substep>

---

<substep id="3d">
<purpose>Proposer/Challenger/Arbiter debate — up to 3 rounds. Each role is a separate Agent call. Never collapse multiple roles into one agent call — self-critique by a single model defeats the purpose.</purpose>

<role id="proposer" model="sonnet">
Synthesizes Subagent A + B output (and `ci_run_evidence` if present); proposes the most likely root cause and a concrete fix with file and line reference.

- If Subagent A returned `confidence: "high"` or `"low"`: inject filtered facts as seed evidence — anchor the hypothesis in the code structure from Subagent B, treating `fix-flaky-test` facts as supporting evidence, not the conclusion.
- If Subagent A returned `confidence: "none"`: note this explicitly and rely solely on Subagent B (and `ci_run_evidence` if present).
- **If `ci_run_evidence` is non-null**: include it in the prompt under a clearly labeled section: *"CI run forensic evidence (single run, not aggregated — weigh accordingly):"*. This data is unscored and reflects only one specific failure, so corroborating signals across `trunk_filtered_facts` and code analysis should outweigh isolated CI-run observations when they conflict. If `ci_run_evidence.status == "build_failure"`, note that test-level data is unavailable from this source.
- Must explicitly state any approaches excluded due to previous failed attempts.
- If any `previous_attempts` entry has a non-null `recommended_next_step`, prepend to prompt: *"Prior investigation ({date}, {outcome}) recommended: '{recommended_next_step}'. Approaches already tried and rejected: {excluded_approaches}. Rejected because: {rejection_reasons}. Start from this hypothesis — confirm or refute it with code evidence before proposing anything else."*
- If 3b-ii returned `classification = "SUT"` with user override (option b or c), prepend: *"Note: this was originally classified SUT (signals: {sut_signals_matched}). The SUT hypothesis: {sut_description}. The test fix should defensively address this."*

<output-schema>
```json
{ "root_cause": "string", "fix_file": "string", "fix_line": "integer", "fix_description": "string", "excluded_approaches": ["string"] }
```
</output-schema>
</role>

<role id="challenger" model="opus">
Receives the Proposer's full output. Challenges the proposal — alternative causes, edge cases, risk of breaking other tests. Must explicitly take a position on whether the proposed causal mechanism itself is sound.

<output-schema>
```json
{ "challenges": ["string"], "mechanism_rebutted": "boolean" }
```
`challenges` must contain at least one item. `mechanism_rebutted: true` means the causal chain was explicitly challenged (e.g. "the proposed collision cannot occur because each run deploys its own contract"), not merely that alternatives were proposed.
</output-schema>
</role>

<role id="arbiter" model="opus">
Receives both Proposer and Challenger outputs. Decides whether to stop (enough confidence) or run another round (max 3 total). Issues the final verdict.

<output-schema>
```json
{ "verdict": "PROCEED | INCONCLUSIVE", "rationale": "string", "next_round": "boolean" }
```
</output-schema>
</role>

<schema-validation applies-to="proposer challenger arbiter">
Same two classes as 3b (transient → immediate retry; structural → retry with error context). Allow up to **3 total attempts** per role. After 3 failures → hard stop for this issue. Apply mid-flight abandonment rule. Write Investigation Update comment (OUTCOME = ABANDONED): state which debate role failed and include the validation error; recommended next step: re-run and include last raw output verbatim. Continue with other issues.
</schema-validation>

</substep>

---

<per-issue-return-schema>
Never surface raw Proposer/Challenger/Arbiter responses to the top-level parent. Distill and return only:

```json
{
  "key": "KEY-NNN",
  "outcome": "PROCEED | INCONCLUSIVE | MISMATCH | SKIP_TOP_LEVEL | RETURNED_TO_QUEUE | ABANDONED",
  "trunk_investigation_status": "existing | triggered | uninvestigated | ci_run_only",
  "trunk_fact_count": "integer",
  "trunk_analysis_url": "string | null",
  "trunk_test_case_url": "string | null",

  "fix_file": "string (PROCEED only, null otherwise)",
  "fix_line": "integer (PROCEED only, null otherwise)",
  "fix_description": "string (PROCEED only, null otherwise)",
  "proposer_root_cause": "string (PROCEED or INCONCLUSIVE only, null otherwise)",
  "excluded_approaches": ["string (PROCEED only, null otherwise)"],
  "classifier": {
    "classification": "TEST | SUT | AMBIGUOUS | INFRA",
    "sut_score": "number", "test_score": "number",
    "sut_signals_matched": ["string"], "test_signals_matched": ["string"],
    "sut_description": "string | null",
    "sut_pivot": { "file": "string | null", "component": "string | null", "hypothesis": "string | null" }
  },

  "fix_description_attempted": "string (INCONCLUSIVE only, null otherwise)",
  "inconclusive_reason": "string (INCONCLUSIVE only, null otherwise)",
  "recommended_next_step": "string | null (INCONCLUSIVE only)",

  "mismatch_details": "string (MISMATCH only, null otherwise)"
}
```

Outcomes RETURNED_TO_QUEUE, ABANDONED, CLOSED_SUBTEST, and SKIP_TOP_LEVEL are **fully handled within the per-issue subagent** (JIRA comment written, abandonment rule applied) before returning. The parent only records the outcome.
</per-issue-return-schema>

---

<checkpoint model="haiku">
<purpose>Print summary and wait for user confirmation before any files are modified. Skip in `--auto` mode — proceed with all PROCEED verdicts automatically (MISMATCH issues were already resolved above).</purpose>

<summary-table>
| Issue | Trunk | Trunk link | Proposed fix location | Verdict |
|-------|-------|------------|-----------------------|---------|
| KEY-123 | existing (2 facts ≥0.9) / triggered (0 facts ≥0.9) / uninvestigated | [Analysis]({trunk_analysis_url}) or [Test case]({trunk_test_case_url}) | `pkg/foo/bar_test.go:447` | PROCEED / INCONCLUSIVE / SKIP_TOP_LEVEL / MISMATCH |

- Use `trunk_analysis_url` for the link when available; fall back to `trunk_test_case_url`.
- `uninvestigated` — Trunk returned no results within 5 minutes; fix based on code analysis only.
- `0 facts ≥0.9` — investigation existed but all facts were below threshold; treated as code-analysis-only.
- `MISMATCH` — the innermost failing function from the Trunk stack trace no longer exists in the codebase.
</summary-table>

<mismatch-handling>
For each MISMATCH issue, show mismatch details inline and ask explicitly — **never auto-resolve in `--auto` mode**:

<user-prompt>
KEY-NNN: The Trunk stack trace references `{function}` in `{file}` which no longer exists in the codebase. The failure data may be outdated. How would you like to proceed?
(a) Investigate anyway using code analysis only (treat as if no Trunk data)
(b) Skip and return ticket to queue
(c) Update the Trunk ticket and retry later
</user-prompt>

Apply the user's choice:
- (a) → re-run this issue through 3d with `trunk_investigation_status = "uninvestigated"`.
- (b) or (c) → apply mid-flight abandonment rule immediately.
</mismatch-handling>

State explicitly: "Investigation is done. Here's the summary above." Then ask: "Proceed with fixes? Exclude specific issues by listing their keys."

If the user excludes or skips any ticket, apply the mid-flight abandonment rule to it immediately.
</checkpoint>

<on_complete>
Write investigation results into `ticket_records` (update `actionable_facts`, `chosen_fix`, outcome fields per ticket).

- Any PROCEED verdicts exist → Read [phase4-apply-fix.md](phase4-apply-fix.md) and follow its instructions.
- All verdicts are INCONCLUSIVE / MISMATCH / SKIP_TOP_LEVEL / RETURNED_TO_QUEUE → Read [phase6-jira-update.md](phase6-jira-update.md) and follow its instructions.
</on_complete>

</phase>
