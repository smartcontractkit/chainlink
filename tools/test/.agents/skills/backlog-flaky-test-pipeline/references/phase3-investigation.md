---
phase: phase3
model: sonnet
---

<phase id="phase3">

<purpose>
Run parallel per-ticket investigations: extract Trunk data, analyze test code, classify flakiness source, debate root cause, and return a structured verdict per ticket. Updates `ticket_records` with `actionable_facts`, `chosen_fix`, and outcome. Runs a checkpoint gate before any files are modified.
</purpose>

<prereqs>

**If `invocation.mode = "local"`: skip this section entirely — no JIRA writes occur in local mode.**

Otherwise (JIRA modes only):
- Read [../../_shared-jira-flaky-ops/investigation-comment.md](../../_shared-jira-flaky-ops/investigation-comment.md) if not already loaded — required for Investigation Update comment format and parsing rules.
- Read [../../_shared-jira-flaky-ops/abandon-ticket.md](../../_shared-jira-flaky-ops/abandon-ticket.md) if not already loaded — required for mid-flight abandonment procedure.
</prereqs>

<parallelism>
Spawn all N per-ticket investigation subagents in a **single message**. Each receives: the slim record (`jira_key`, `local_id`, `title`, `description`, `trunk_test_case_url`, `test_name`, `package`, `previous_attempts`, `provided_log_text`), `mode` (`project` | `direct-ticket` | `local`), `nav_tool`, `lsp_available`, `auto_mode`. Never pass raw JIRA API objects.
</parallelism>

---

<substep id="3a" model="haiku">
<purpose>Resolve `testCaseId`, retrieve `fix-flaky-test` historical data, optionally invoke `investigate-ci-failure` for single-run forensic data, and run the top-level subtest check.</purpose>

<local-mode-branch>
**If `mode = "local"`**: skip the `fix-flaky-test-call` and `investigate-ci-failure-call` blocks entirely.

- If `slim_record.provided_log_text` is non-null → set `trunk_filtered_facts = [slim_record.provided_log_text]`, `trunk_investigation_status = "user_provided"`.
- If `provided_log_text` is null → set `trunk_filtered_facts = []`, `trunk_investigation_status = "uninvestigated"`.

Proceed directly to `<top-level-check>`.
</local-mode-branch>

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
After Trunk investigation resolves (or after the local-mode branch sets `trunk_filtered_facts`), inspect `trunk_filtered_facts` (and any stack trace within them) for a `file:line`. If the failure line falls inside a `t.Run(...)` callback AND the outer function contains no assertions outside `t.Run` blocks → candidate for **SKIP_TOP_LEVEL**.

**Exception**: if `slim_record.test_name` contains a `/` the ticket was already filed against a specific subtest — SKIP_TOP_LEVEL must not fire. If `test_name` is null, check the title for `/`.

If all conditions are met:

**JIRA mode** (mode ≠ "local"):
1. Follow `_shared-jira-flaky-ops/transition-ticket.md` with `jira_key` and `target = "Won't Do"`.
2. Follow `_shared-jira-flaky-ops/investigation-comment.md` to write `addCommentToJiraIssue` (OUTCOME = CLOSED_SUBTEST). "What was investigated": failure originates in a `t.Run` subtest, not the top-level function. "Recommended next step": file or locate a ticket for the specific subtest. All other sections: N/A.
3. Stop investigation for this issue.

**Local mode** (mode = "local"):
- Print: *"Skipping {test_name} — failure originates in a t.Run subtest, not the top-level function. File a ticket against the specific subtest if you want this fixed."*
- Return verdict `SKIP_TOP_LEVEL` for this record. No JIRA writes.
</top-level-check>

</substep>

---

<substep id="3b" parallelism="Subagent A and Subagent B spawn in a single message">

<subagent id="trunk-analyzer" model="haiku">
<inputs>
`trunk_filtered_facts` (already filtered to ≥ 0.9 in 3a, or the user-provided log in local mode), `trunk_investigation_status`, `mode`. Do not call any Trunk MCP tools.
</inputs>
<logic>
- If `trunk_investigation_status = "uninvestigated"` or `trunk_filtered_facts` is empty → return `confidence: "none"` with empty `facts`.
- If `trunk_investigation_status = "user_provided"` (local mode) → treat `provided_log_text` as observational evidence with `confidence: "low"`. It is symptom data, not aggregated CI data.
- Map confidence from pre-filtered facts (JIRA modes):
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

Allow up to **3 total attempts** per subagent. After 3 failures → hard stop for this issue. In **JIRA mode**: follow `_shared-jira-flaky-ops/abandon-ticket.md` then write Investigation Update comment via `_shared-jira-flaky-ops/investigation-comment.md` (OUTCOME = ABANDONED): state which subagent failed and include the validation error; recommended next step: re-run and include last raw output verbatim. In **local mode**: no JIRA writes — just return verdict `ABANDONED` with the error. Continue with other issues.
</schema-validation>

</substep>

---

<substep id="3b-ii" model="sonnet" runs-in="parent">
<purpose>
Classify flakiness source as TEST / SUT / INFRA / AMBIGUOUS before entering the fix debate. Runs in the parent (not a subagent) because it may require a user gate. Single LLM call — no scoring, no signal enumeration, no tier ladder. The model examines the available evidence and chooses one classification, with every conclusion grounded in a verbatim quote.
</purpose>

<input-schema>
```json
{
  "$schema": "phase_3bii_input_v1",
  "record_id": "string (jira_key or local_id)",
  "test_name": "string",
  "trunk_filtered_facts": ["string"],
  "trunk_investigation_status": "existing | triggered | uninvestigated | ci_run_only | user_provided",
  "subagent_b_output": {
    "file": "string", "line": "number", "analysis": "string",
    "suspected_cause": "string", "suspected_cause_location": "test_code | production_code | unknown",
    "excluded_approaches": ["string"], "code_mismatch": "boolean", "mismatch_details": "string | null"
  }
}
```
</input-schema>

<classifier-call model="sonnet">
A single Sonnet call. The prompt:

> You are classifying a flaky test failure. Choose **one** of:
> - **TEST** — the test code introduces non-determinism (timing dependency, shared state, missing cleanup, parallelism without sync, non-deterministic data, ordering assumption, race in test code, hardcoded resources, etc.).
> - **SUT** — the production code under test is incorrect, racy, or not ready when the test exercises it.
> - **INFRA** — an environmental failure unrelated to either the test or the SUT (OOM, disk full, image pull failure, network outage at infra layer).
> - **AMBIGUOUS** — evidence is insufficient, contradictory, or absent.
>
> **Rules:**
> 1. Every classification must be backed by 1–3 verbatim quotes. A quote's `source` is one of:
>    - `trunk_facts` — excerpt must appear verbatim in `trunk_filtered_facts`.
>    - `code_analysis` — excerpt must appear verbatim in `subagent_b_output.analysis`.
>    - `direct_field` — the literal value `"test_code"` or `"production_code"` from `subagent_b_output.suspected_cause_location` (counts as TEST or SUT evidence respectively).
> 2. **Never quote raw test source code, function names, or code comments.** Only the analyzer's *synthesized* `analysis` text counts as code-side evidence.
> 3. If you cannot produce at least one valid quote for the winning side → classify AMBIGUOUS with confidence `none`.
> 4. INFRA requires at least one `trunk_facts` quote. Code analysis alone cannot establish INFRA.
> 5. If `subagent_b_output.code_mismatch == true` → classify AMBIGUOUS (stale stack trace; attribution unsafe).
> 6. Confidence: `high` = 2+ corroborating quotes, no contradicting evidence; `low` = 1 quote OR thin/indirect quotes; `none` = no usable evidence (only on AMBIGUOUS).
> 7. For SUT: also produce `sut_description` (one sentence) and `sut_pivot` (file/component/hypothesis — fields may be null).
> 8. `pattern_category` is a short free-form label (≤ 5 words) for diagnostic display — e.g. "timing dependency", "OOM during test", "stale precondition", "production nil deref". Not load-bearing.

Inputs to inject into the prompt: `trunk_filtered_facts`, `subagent_b_output` (all fields). Bias the model toward AMBIGUOUS when evidence is thin — false TEST classification leads to bogus fixes; AMBIGUOUS just surfaces the case to the user.
</classifier-call>

<output-schema>
```json
{
  "$schema": "phase_3bii_output_v1",
  "classification": "TEST | SUT | AMBIGUOUS | INFRA",
  "confidence": "high | low | none",
  "rationale": "string (one sentence)",
  "pattern_category": "string | null (≤ 5 words; diagnostic label only)",
  "evidence": [
    {
      "source": "trunk_facts | code_analysis | direct_field",
      "excerpt": "string",
      "supports": "TEST | SUT | INFRA"
    }
  ],
  "sut_description": "string | null",
  "sut_pivot": { "file": "string | null", "component": "string | null", "hypothesis": "string | null" }
}
```
`sut_description` required (non-null) when classification == SUT. `sut_pivot` required (fields may be null) when classification is SUT or AMBIGUOUS; null otherwise.
</output-schema>

<validation>
Post-call, deterministically:

1. **Excerpt verification** — for each `evidence` entry:
   - `trunk_facts` → search the joined `trunk_filtered_facts` text. Excerpt must appear verbatim.
   - `code_analysis` → search `subagent_b_output.analysis`. Excerpt must appear verbatim.
   - `direct_field` → excerpt must equal either `"test_code"` or `"production_code"` AND match `subagent_b_output.suspected_cause_location`.
   Drop any entry that fails verification.
2. **Consistency** — apply in order:
   - `subagent_b_output.code_mismatch == true` → force `classification = AMBIGUOUS`, `confidence = none`.
   - After dropping invalid evidence, if no evidence remains supporting the chosen classification → force AMBIGUOUS.
   - `classification == INFRA` but zero validated `trunk_facts` quotes → force AMBIGUOUS.
   - `classification == SUT` but `sut_description` is null → force AMBIGUOUS.
3. **Schema validation** — transient (empty/null) → retry; structural (missing fields, wrong types, invalid enum) → retry with error context. Up to 3 total attempts. After 3 failures: `classification = AMBIGUOUS`, `confidence = none`, `rationale = "Schema validation failed after 3 attempts"`, `evidence = []`. Continue to gate logic.
</validation>

<gate-logic>
| Classification | `--auto` JIRA mode | Interactive JIRA mode | `--auto` local mode | Interactive local mode |
|---|---|---|---|---|
| TEST | Continue to 3c | Continue to 3c | Continue to 3c | Continue to 3c |
| SUT | Return to queue + JIRA comment | Prompt user (options a/b/c) | Skip test + report | Prompt: (a) skip (b) override to TEST — no JIRA option |
| AMBIGUOUS | Return to queue + JIRA comment | Prompt user with evidence | Skip test + report | Prompt: (a) skip (b) override to TEST |
| INFRA | Return to queue + JIRA comment | Prompt user | Skip test + report | Prompt: (a) skip (b) override to TEST |

In local mode, "skip" means: no JIRA writes, no fix attempted; include in final summary as SKIPPED with the classification reason.

<user-prompt id="sut-gate">
Classification: **SUT** (confidence: {confidence})
Pattern: {pattern_category}
Rationale: {rationale}

Evidence:
{for each evidence row: - "{excerpt}" (from {source}, supports {supports})}

{sut_description}

This test appears to expose a SUT bug, not a test-code bug. Options:
(a) Return this ticket to the queue with a SUT-hypothesis note
(b) Override — treat as TEST and proceed to fix debate (audited)
(c) Fix the test code AND auto-file a SUT bug ticket (label: sut-bug)
</user-prompt>

<user-prompt id="ambiguous-gate">
Classification: **AMBIGUOUS** (confidence: {confidence})
Rationale: {rationale}

Evidence:
{for each evidence row: - "{excerpt}" (from {source}, supports {supports})}
{if evidence is empty: (no quotable evidence — investigation cannot reliably attribute the failure)}

Options:
(a) Return this ticket to the queue / skip
(b) Override — treat as TEST and proceed to fix debate (audited)
</user-prompt>

<user-prompt id="infra-gate">
Classification: **INFRA** (confidence: {confidence})
Pattern: {pattern_category}
Rationale: {rationale}

Evidence:
{for each evidence row: - "{excerpt}" (from {source})}

This test failed for environmental reasons unrelated to the test or SUT. Options:
(a) Return to queue / skip (recommended — re-run when infra is healthy)
(b) Override — treat as TEST and proceed to fix debate (audited)
</user-prompt>

*Option (b) audit trail*: add JIRA comment *"Classification overridden to TEST by user. Original: {classification}, confidence: {confidence}, rationale: {rationale}."* Add commit trailer: `Flakiness-classification: TEST (user override from {classification})`.

*Option (c)* (SUT only): create a JIRA issue in the same project: summary `SUT bug: {sut_description}`, description includes `sut_pivot` fields and the evidence list, label `sut-bug`. Return the new ticket key to the user. Proceed to 3c treating current ticket as TEST.

For SUT/AMBIGUOUS/INFRA auto-queue returns in **JIRA mode**: follow `_shared-jira-flaky-ops/investigation-comment.md` to write `addCommentToJiraIssue` (OUTCOME = RETURNED_TO_QUEUE). "What was investigated": classification, confidence, pattern category, rationale, evidence quotes. "Hypothesis": `sut_description` if SUT, otherwise N/A. "Recommended next step": SUT → investigate `sut_pivot`; AMBIGUOUS → clarify before re-investigating; INFRA → re-run after infra is healthy. Then follow `_shared-jira-flaky-ops/abandon-ticket.md` (unassign + transition to Open). Continue with other issues.

For SUT/AMBIGUOUS/INFRA in **local mode**: no JIRA writes. Return verdict `SKIPPED` with the classification reason. Continue with other issues.
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
- If 3b-ii returned `classification = "SUT"` with user override (option b or c), prepend: *"Note: this was originally classified SUT (rationale: {rationale}; pattern: {pattern_category}). The SUT hypothesis: {sut_description}. The test fix should defensively address this."*

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
Same two classes as 3b (transient → immediate retry; structural → retry with error context). Allow up to **3 total attempts** per role. After 3 failures → hard stop for this issue. In **JIRA mode**: follow `_shared-jira-flaky-ops/abandon-ticket.md` then write Investigation Update comment via `_shared-jira-flaky-ops/investigation-comment.md` (OUTCOME = ABANDONED): state which debate role failed and include the validation error; recommended next step: re-run and include last raw output verbatim. In **local mode**: no JIRA writes — return verdict `ABANDONED` with the error. Continue with other issues.
</schema-validation>

</substep>

---

<per-issue-return-schema>
Never surface raw Proposer/Challenger/Arbiter responses to the top-level parent. Distill and return only:

```json
{
  "key": "KEY-NNN | null",
  "local_id": "local-N | null",
  "record_id": "KEY-NNN (jira_key if non-null, else local_id)",
  "outcome": "PROCEED | INCONCLUSIVE | MISMATCH | SKIP_TOP_LEVEL | RETURNED_TO_QUEUE | ABANDONED | SKIPPED",
  "subagent_calls_made": ["trunk_analyzer | code_analyzer | classifier | proposer | challenger | arbiter"],
  "trunk_investigation_status": "existing | triggered | uninvestigated | ci_run_only | user_provided",
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
    "confidence": "high | low | none",
    "rationale": "string",
    "pattern_category": "string | null",
    "evidence": [{ "source": "trunk_facts | code_analysis | direct_field", "excerpt": "string", "supports": "TEST | SUT | INFRA" }],
    "sut_description": "string | null",
    "sut_pivot": { "file": "string | null", "component": "string | null", "hypothesis": "string | null" }
  },

  "fix_description_attempted": "string (INCONCLUSIVE only, null otherwise)",
  "inconclusive_reason": "string (INCONCLUSIVE only, null otherwise)",
  "recommended_next_step": "string | null (INCONCLUSIVE only)",

  "mismatch_details": "string (MISMATCH only, null otherwise)"
}
```

Outcomes RETURNED_TO_QUEUE, ABANDONED, CLOSED_SUBTEST, and SKIP_TOP_LEVEL are **fully handled within the per-issue subagent** (JIRA comment written and abandonment rule applied in JIRA modes; no JIRA writes in local mode) before returning. The parent only records the outcome.

`SKIPPED` is the local-mode equivalent of `RETURNED_TO_QUEUE` — no JIRA writes, included in the final summary with the classification reason.

</per-issue-return-schema>

---

<integrity-gate>
**The parent MUST check `subagent_calls_made` on every per-issue return. This check is never skipped — not in `--auto` mode, not for tickets that look obviously simple, not when the verdict appears self-evidently correct. The check is the only thing standing between the multi-agent protocol and a single model rationalizing its way to a confidently-wrong conclusion.**

Required entries by outcome:

| Outcome | Required entries in `subagent_calls_made` |
|---------|-------------------------------------------|
| `PROCEED` / `INCONCLUSIVE` | all six: `trunk_analyzer`, `code_analyzer`, `classifier`, `proposer`, `challenger`, `arbiter` |
| `RETURNED_TO_QUEUE` / `SKIPPED` | `trunk_analyzer`, `code_analyzer`, `classifier` (classification ran; debate did not) |
| `MISMATCH` | `trunk_analyzer`, `code_analyzer` (short-circuits at 3c) |
| `SKIP_TOP_LEVEL` | `trunk_analyzer` (short-circuits in 3a) |
| `ABANDONED` | whatever ran before the abandonment trigger; document the trigger in the return |

If any required entry is missing → **override the outcome to `ABANDONED`**, reason `"protocol skipped — missing: {list of roles}"`. Do not apply the fix. Do not negotiate. Do not accept the verdict on the grounds that the inline reasoning "looks right."
</integrity-gate>

---

<checkpoint model="haiku">
<purpose>Print summary and wait for user confirmation before any files are modified. Skip in `--auto` mode — proceed with all PROCEED verdicts automatically (MISMATCH issues were already resolved above).</purpose>

<summary-table>
**JIRA modes** — include Trunk columns:

| Issue | Trunk | Trunk link | Proposed fix location | Verdict |
|-------|-------|------------|-----------------------|---------|
| KEY-123 | existing (2 facts ≥0.9) / triggered (0 facts ≥0.9) / uninvestigated | [Analysis]({trunk_analysis_url}) or [Test case]({trunk_test_case_url}) | `pkg/foo/bar_test.go:447` | PROCEED / INCONCLUSIVE / SKIP_TOP_LEVEL / MISMATCH |

- Use `trunk_analysis_url` for the link when available; fall back to `trunk_test_case_url`.
- `uninvestigated` — Trunk returned no results within 5 minutes; fix based on code analysis only.
- `0 facts ≥0.9` — investigation existed but all facts were below threshold; treated as code-analysis-only.
- `MISMATCH` — the innermost failing function from the Trunk stack trace no longer exists in the codebase.

**Local mode** — drop Trunk columns; use `local_id + test_name` in the Issue column:

| Issue / Test | Evidence | Proposed fix location | Verdict |
|--------------|----------|----------------------|---------|
| local-1 · TestFoo | user log / none | `pkg/foo/bar_test.go:447` | PROCEED / INCONCLUSIVE / SKIPPED / MISMATCH |

- `Evidence` = "user log" if `provided_log_text` non-null, else "none".
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

If the user excludes or skips any ticket: in JIRA mode, follow `_shared-jira-flaky-ops/abandon-ticket.md` for it immediately. In local mode, no JIRA writes — just exclude from fix list.
</checkpoint>

<on_complete>
Write investigation results into `ticket_records` (update `actionable_facts`, `chosen_fix`, outcome fields per ticket).

- Any PROCEED verdicts exist → Read [phase4-apply-fix.md](phase4-apply-fix.md) and follow its instructions.
- All verdicts are INCONCLUSIVE / MISMATCH / SKIP_TOP_LEVEL / RETURNED_TO_QUEUE / SKIPPED (no PROCEED) and `mode = "local"` → Read [phase-final-local.md](phase-final-local.md) and follow its instructions.
- All verdicts are INCONCLUSIVE / MISMATCH / SKIP_TOP_LEVEL / RETURNED_TO_QUEUE (no PROCEED) and `mode ≠ "local"` → Read [phase6-jira-update.md](phase6-jira-update.md) and follow its instructions.
</on_complete>

</phase>
