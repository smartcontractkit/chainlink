---
phase: phase2a
model: haiku
---

<phase id="phase2a">

<purpose>
Fetch and filter JIRA issues for a project key, resolve each to a local test function, and build slim ticket records. Populates `ticket_records` in context.
</purpose>

<prereqs>

Read [../../_shared-jira-flaky-ops/investigation-comment.md](../../_shared-jira-flaky-ops/investigation-comment.md) before parsing any previous-investigation comments. It defines the comment format and parsing rules.
</prereqs>

<parent-setup>
Run `git remote get-url origin` once. Extract `{owner}/{repo}` from the URL. Cache as `current_repo`. Pass to the subagent below.
</parent-setup>

<subagent id="fetch-filter" model="haiku">

Spawn this subagent following `_shared-jira-flaky-ops/fetch-flaky-tickets.md`. Read that file for the full loop, filtering rules, and field extraction logic.

Inputs to pass: `KEY`, `N`, `cloudId`, `current_repo`, `nav_tool`, `lsp_available` (from `phase_outputs.phase0`), repo root path.

Expected output:
```json
{
  "slim_records": [ /* see _shared-jira-flaky-ops/README.md for schema */ ],
  "skipped": { "cross_repo": 0, "system_tests": 0, "not_found": 0 }
}
```

</subagent>

<on_subagent_return>
Write slim records to `ticket_records` in context.

If `len(slim_records) < N`: inform user:
> Found K eligible issues. Skipped: {cross_repo} cross-repo, {system_tests} system-tests (excluded), {not_found} test function not found locally.
</on_subagent_return>

<on_complete>
Read [phase2c-prior-gate.md](phase2c-prior-gate.md) and follow its instructions.
</on_complete>

</phase>
