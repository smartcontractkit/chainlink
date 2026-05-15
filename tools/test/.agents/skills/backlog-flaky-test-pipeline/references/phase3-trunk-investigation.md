---
name: phase3-trunk-investigation
description: Shared procedure — retrieve Trunk and CI-run evidence for a JIRA-mode ticket. Covers fix-flaky-test (Branch A/B/C) and investigate-ci-failure with recovery. Called from phase3a-ii-jira-evidence.md.
---

<procedure id="trunk-investigation">

<purpose>
Retrieve observational evidence from Trunk and optionally from a specific CI run. Writes `actionable_facts`, `ci_run_evidence`, `trunk_facts_quality`, and `trunk_analysis_url` to the calling context.
</purpose>

<note>
Only called in JIRA modes (`mode != "local"`). Apply the `Confidence >= 0.9` filter immediately on every `fix-flaky-test` response — never after-the-fact.
</note>

<quality-classification>
After every `fix-flaky-test` response, once `actionable_facts` is populated (or confirmed empty), set `trunk_facts_quality`:

- `actionable_facts` is empty → `trunk_facts_quality = "none"`.
- At least one fact contains raw CI observational data — exact log lines, error messages, stack traces, or specific `file:line` from failing runs → `trunk_facts_quality = "aggregated"`.
- Facts are present but describe symptoms only (descriptive paraphrasing, e.g. *"failures in cluster 2 are regex mismatches"*) with no raw data → `trunk_facts_quality = "symptom_only"`.

This classification is final once written — do not re-evaluate downstream.
</quality-classification>

<inputs>
- `test_case_id` — UUID from `customfield_13010`; null if not available.
- `test_name` — for fuzzy search fallback.
- `title` — for fuzzy search fallback.
- `ci_run_url` — GitHub Actions run URL; null if not provided upfront.
- `jira_key` — for user-facing messages.
- `auto_mode` — boolean.
</inputs>

<step id="fix-flaky-test">
Use `test_case_id`. If null: call `mcp__trunk__search-test` with `test_name` (falling back to `title`) — flag if used, it may match the wrong test.

Call `mcp__trunk__fix-flaky-test` with `testCaseId` (no `createNewInvestigation`). **Immediately after each response (including polls), apply the filter:**
1. Capture `trunk_analysis_url` from `investigationUrl` if present.
2. Extract only facts with `Confidence >= 0.9` from `## Facts`. Discard lower-confidence facts and the entire `## Markdown Summary`.
3. Store as `actionable_facts: [str]`.

**If response has `summary` or `facts`**: classify per `<quality-classification>`. If `ci_run_url` is non-null, proceed to `<step id="investigate-ci-failure">`. Otherwise done.

**If response is empty** (no investigation yet) — branch on `ci_run_url`:

**Branch A: `ci_run_url` is non-null** → skip polling; rely on `investigate-ci-failure`. `actionable_facts` stays `[]` → `trunk_facts_quality = "none"`. Proceed to `<step id="investigate-ci-failure">`.

**Branch B: `ci_run_url` is null AND `auto_mode` is false** → ask using `AskUserQuestion`:
- Question: "{jira_key}: No existing Trunk investigation found. How would you like to proceed?"
- Options:
  - "Attach a CI run URL" — user provides URL via the Other field → set `ci_run_url`, proceed to `investigate-ci-failure`; `trunk_facts_quality = "none"`.
  - "Trigger a new investigation (2–5 min)" → fall through to Branch C.
  - "Skip this ticket" → `trunk_facts_quality = "none"`; return.

**Branch C: no CI URL and (`auto_mode` or user chose "trigger")** → trigger:
- Inform user: *"No existing Trunk investigation found for {jira_key}. Triggering one now — this may take 2–5 minutes."*
- Call `mcp__trunk__fix-flaky-test` with `createNewInvestigation: true`.
- Poll every ~30 seconds for up to 5 minutes (10 attempts). Apply the same immediate filter on each response.
- If `investigationId` is unknown during polling: expected — still initializing. Inform user: *"Investigation pending (investigationId not yet active — this is normal). Continuing to poll…"*
- Non-empty result → classify per `<quality-classification>`. Done.
- Still empty after 5 minutes → `trunk_facts_quality = "none"`. Ask using `AskUserQuestion`:
  - Question: "Trunk did not return investigation results after 5 minutes for {jira_key}. How would you like to proceed?"
  - Options: ["Continue with code analysis only", "Wait longer", "Skip this ticket"]
</step>

<step id="investigate-ci-failure">
Runs when `ci_run_url` is non-null.

Call `mcp__trunk__investigate-ci-failure` with `workflowUrl = ci_run_url` and `orgSlug = "chainlink"` (hardcoded — see `trunk-org-slug` tip in SKILL.md).
Store the full structured response as `ci_run_evidence`. **No `Confidence >= 0.9` filter — separate evidence track per the `<absolute_constraints>` exemption.**

Edge cases:
- "Build/compile failure before tests ran" → `ci_run_evidence = { "status": "build_failure", "raw": <response> }`. Inform user: *"CI run for {jira_key} failed before tests executed; CI evidence will be limited."* Continue.
- Tool errors (permission, unknown error, malformed URL) → `ci_run_evidence = null`, log the error. **If `actionable_facts` is also empty** (zero Trunk evidence — typical of Branch A entry), recover:
  - `auto_mode = true` → trigger a new investigation without prompting. Inform user: *"investigate-ci-failure failed for {jira_key} ({error}); auto mode falling back to fix-flaky-test (2–5 minutes)."* Run Branch C.
  - `auto_mode = false` → ask using `AskUserQuestion`:
    - Question: "`investigate-ci-failure` failed for {jira_key} ({error}). Trunk has no investigation data. How would you like to proceed?"
    - Options: ["Trigger a new fix-flaky-test investigation (2–5 min)", "Continue with code analysis only", "Skip this ticket"]
    - "Trigger" → run Branch C.
    - Other → leave `actionable_facts = []`, continue (diagnose-fallback in the caller may still fire).
</step>

<output>
Side-effects written to the calling context:
- `actionable_facts: [str]` — facts with `Confidence >= 0.9`, or `[]`
- `ci_run_evidence: object | null`
- `trunk_facts_quality: "aggregated" | "symptom_only" | "none"`
- `trunk_analysis_url: string | null`

Note: `ci_run_evidence` is not consumed by the 3b-ii classifier — only by the Proposer in substep 3d.
</output>

</procedure>
