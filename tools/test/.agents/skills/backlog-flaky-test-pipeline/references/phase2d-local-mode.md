---
phase: phase2d
model_tier: lightweight
---

<phase id="phase2d">

<purpose>
Build slim records for local mode (no JIRA, no Trunk). One record per test spec provided via `--local`. Replaces phases 2a/2b/2c for local invocations. Populates `ticket_records` in context.
</purpose>

<note>
This phase only runs when `invocation.mode = "local"`. It never touches JIRA or Trunk. No tickets are claimed. No prior-attempt gating applies.
</note>

<steps>

<step id="1-parse-specs">
Parse the test specs from `invocation.args.test_specs` (populated by phase1). Each spec has the form `<package>.<TestName>` (period separator) or bare `TestName` (no period). Subtest segments use `/`.

For each spec:
- If a period is present: left side = package import path, right side = test function name (including any `/subtest` suffix).
- If no period: treat the whole token as `TestName`; package will be resolved via code-nav in the next step.
</step>

<step id="2-locate-test">
For each test spec, locate the test function in the repo:

1. **LSP definition lookup** (if `nav_tool = "lsp"` or `lsp_available = true`): look up the definition of `func {TestName}`.
2. **code-review-graph** (if `nav_tool = "crg"`): `mcp__code-review-graph__semantic_search_nodes_tool` with the test name.
3. **Last resort**: `grep -r "func {TestName}" .` — parse first `filepath:line`, warn if multiple matches.

If not found after all three: report *"Test function {TestName} not found in the repo — skipping."* and do not create a record for this spec. Continue with remaining specs.

If the package was not provided in the spec, infer it from the found file path (convert file path to Go import path using the module path from `go.mod`).
</step>

<step id="3-read-log">
If `invocation.args.log_path` is non-null and has not yet been read:
- Read the file once into a single string `provided_log_text`.
- If the file is missing: warn user *"--log file not found at {path}; proceeding without log evidence."* and set `provided_log_text = null`.

Apply to all records (one log file shared across all test specs).
</step>

<step id="4-build-records">
For each located test, build a slim record conforming to the schema in `_shared-jira-flaky-ops/README.md`:

```json
{
  "jira_key":            null,
  "local_id":            "local-{N}",
  "title":               "{TestName}",
  "description":         "",
  "trunk_test_case_url": null,
  "test_case_id":        null,
  "package":             "<resolved package import path>",
  "test_name":           "<TestName or TestName/subtest>",
  "previous_attempts":   [],
  "ci_run_url":          null,
  "provided_log_path":   "<invocation.args.log_path or null>",
  "provided_log_text":   "<string or null>"
}
```

Assign `local_id` values sequentially starting at `local-1`.
</step>

<step id="5-announce">
Write all built records to `ticket_records` in context.

Announce: *"Local mode: {N} tests prepared. No JIRA tickets claimed. Proceeding to investigation."*

If N = 0 (all specs failed to locate): stop. Nothing to investigate.
</step>

</steps>

<on_complete>
Read [phase3-investigation.md](phase3-investigation.md) and follow its instructions.
</on_complete>

</phase>
