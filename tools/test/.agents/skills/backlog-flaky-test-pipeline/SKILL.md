---
name: backlog-flaky-test-pipeline
description: >-
  A high-automation, multi-agent workflow for resolving flaky tests. Supports
  three modes: JIRA project backlog, direct JIRA ticket IDs, and local benchmark
  mode (no JIRA or Trunk required).

  USE THIS WHEN:
  1. You have specific JIRA ticket IDs (e.g., CRE-123) or want to pull N tickets from a project backlog.
  2. The failure is identified as "flaky" by Trunk.io analysis.
  3. You need an end-to-end pipeline (fetch -> debate -> fix -> PR -> JIRA update).
  4. You want a one-shot multi-agent fix attempt on named tests without JIRA (--local benchmark mode).

  DO NOT USE THIS WHEN:
  1. Iterative, interactive debugging of a single test file — use 'debug-test-failure' instead.
  2. You need an exploratory diagnostic loop rather than a structured fix attempt.
---

<overview>
End-to-end workflow: fetch flaky-test JIRA tickets → run Trunk + code analysis in parallel → apply fixes → commit → update JIRA with the PR. In local mode: take test names directly, skip JIRA + Trunk, stop after fix verification.
</overview>

<usage>
- `/backlog-flaky-test-pipeline [KEY [N]] [--auto]` — project + count mode
- `/backlog-flaky-test-pipeline PROJ-NNN[@<run-url>] [PROJ-NNN[@<run-url>] ...] [--auto]` — direct-ticket mode (e.g. `CRE-5719 CCIP-42@https://github.com/.../actions/runs/123`). The optional `@<run-url>` attaches a CI run URL for `investigate-ci-failure` analysis (direct-ticket mode only).
- `/backlog-flaky-test-pipeline --local <pkg>.TestName [<pkg>.TestName/subtest ...] [--log <path>] [--auto]` — local benchmark mode. Takes test specs directly; no JIRA or Trunk required. Stops after phase4 verification (no commit, no PR). Use for benchmarking this skill against other tools.
- `KEY` — JIRA project key (e.g. `CRE`). Skips Phase 1 prompt if provided.
- `N` — number of issues (default 3). Skips Phase 1 prompt if provided.
- `--auto` — accept all defaults at every gate; only blocks on hard failures and ownership conflicts.
- `--log <path>` — (local mode only) path to a log file whose contents become seed evidence for analysis.
</usage>

<model-assignment>
Set the `model` parameter explicitly when invoking Agent tool calls:
- `[haiku]` — pure API calls or mechanical branching; no reasoning needed
- `[sonnet]` — code understanding or moderate judgment
- `[opus]` — adversarial reasoning or high-stakes go/no-go decisions
</model-assignment>

<absolute_constraints>
- **User gates belong in the parent, never in subagents.** Subagents are fire-and-return and cannot pause for decisions. Resolve all user-facing choices before spawning subagents; they receive already-resolved inputs.
- **Only act on `fix-flaky-test` facts with `Confidence >= 0.9`.** Discard lower-confidence facts entirely. Never use the `## Markdown Summary` — it blends reliable and speculative inferences into a narrative that anchors analysis on wrong causes. (Exemption: `investigate-ci-failure` output is single-run forensic data without a confidence field — kept on a separate evidence track, see `ci_run_evidence` below.)
- **Check go.mod before writing any new utility code.** Three lines of existing library usage beats 30 lines of hand-rolled logic that has to be maintained and tested.
- **Never alter a test's core goal to make it pass.** Changing what a test asserts to eliminate a failure is not a fix.
- **Never remove tests or assertions** unless replacing them with stronger coverage or deleting confirmed dead code.
- **Never modify package-wide helpers to fix a single test.** Scope fixes to the test or the code under test, not shared infrastructure.
</absolute_constraints>

<context_compaction>
This skill is a multi-phase workflow; later phases depend on outputs from earlier ones. The records below must persist across compaction as structured fields — never collapse them into prose. This schema is **authoritative**: if a phase needs a new field, add it here first, do not introduce ad-hoc state.

<never_compact>
  <skill_rules>
    The full `<absolute_constraints>` and `<model-assignment>` blocks from this file. They are referenced in every phase.
  </skill_rules>

  <invocation>
    - `mode`: `project` | `direct-ticket` | `local`
    - `args`: original user input verbatim (project key, count, ticket list, flags). In local mode: `{ test_specs: ["<pkg>.<TestName>", ...], log_path: "<path> | null" }`.
    - `auto_mode`: boolean
  </invocation>

  <phase_state>
    - `current_phase`: e.g. `phase3`
    - `completed_phases`: ordered list
    - `phase_outputs`: keyed by phase id, so a later phase can re-read what an earlier one produced
  </phase_state>

  <ticket_records>
    One row per ticket. Maintain as a table, not a paragraph. Required fields:
    - `jira_key` (e.g. `CRE-5719`) — null in local mode
    - `local_id` (e.g. `local-1`) — null in JIRA modes
    - `test_case_id` (UUID from `customfield_13010`) — null in local mode
    - `test_name`, `package_path`
    - `trunk_investigation_status` — e.g. `"existing" | "triggered" | "uninvestigated" | "ci_run_only" | "user_provided" | "diagnose_run"`; null in local mode before 3a runs
    - `actionable_facts` — facts with `Confidence >= 0.9` from Trunk, or raw strings from `diagnose`/user log
    - `ci_run_url` — GitHub Actions run URL for `investigate-ci-failure` (null in project mode and local mode; populated from `KEY@URL` syntax or 3a fallback prompt in direct-ticket mode)
    - `ci_run_evidence` — structured failure data from `investigate-ci-failure` (null if no URL provided or call failed). Separate evidence track — exempt from the ≥ 0.9 rule. Always null in local mode.
    - `provided_log_path` — path supplied via `--log` in local mode; null in JIRA modes
    - `provided_log_text` — contents of the log file read once; null if no log or file missing
    - `shared_cause_group` — id linking tickets in the same package that share a root cause (see `shared-package-cause` tip)
    - `chosen_fix` — { fix_file, fix_line, fix_description, proposer_root_cause }
    - `applied_fix` — { diff }
    - `lint_status` — `"ran"` | `"skipped"` | `"failed"`
    - `lint_scope` — package path used for linting (e.g. `./core/chains/evm/txmgr/...`)
    - `verification` — { iterations_passed, ci_status }
    - `branch`, `commit_sha`, `pr_url` — null in local mode
  </ticket_records>

  <user_decisions>
    Gate-by-gate log: { `gate_name`, `user_choice`, `rationale` }. Decisions are point-in-time — once recorded, never re-ask.
  </user_decisions>
</never_compact>

<safe_to_compact>
- Verbose tool output once its actionable facts are extracted into a ticket record.
- Subagent transcripts once their conclusions are recorded in the relevant ticket record or `user_decisions` log.
- Intermediate reasoning chains once a decision is committed.
</safe_to_compact>
</context_compaction>

<tips>
<tip id="test-case-id">`testCaseId` comes from `customfield_13010` (TrunkID) — a bare UUID, no URL parsing needed. If absent, fall back to extracting the UUID from the Trunk URL in the description. Use `search-test` fuzzy fallback only as a last resort — it may match the wrong test if names are similar.</tip>

<tip id="mcp-suffix">All code-review-graph MCP tool names carry a `_tool` suffix (e.g. `get_minimal_context_tool`, `semantic_search_nodes_tool`). If these tools do not appear in your environment, skip the crg path and require LSP.</tip>

<tip id="quarantine">Quarantined tests still run in CI (`RUN_QUARANTINED_TESTS=true`), so a successful local 10x run is a reasonable signal but not definitive.</tip>

<tip id="trunk-polling">When Trunk returns that an `investigationId` is unknown during polling, this is expected — it means the investigation is still initializing, not that it failed.</tip>

<tip id="shared-package-cause">When multiple tickets target tests in the same package, check whether they share a common failure cause before analyzing each independently. A single root cause can surface as several distinct JIRA tickets.</tip>

<tip id="trunk-org-slug">When calling `mcp__trunk__investigate-ci-failure`, always pass `orgSlug = "chainlink"`. The Trunk org slug is independent of the GitHub org and the call returns a generic "Unknown error" if it's omitted (account is a member of multiple Trunk orgs).</tip>
</tips>

<begin>

Read [phase0-prerequisites.md](./references/phase0-prerequisites.md) and follow its instructions.
</begin>
