---
phase: phase2a
model: haiku
---

<phase id="phase2a">

<purpose>
Fetch and filter JIRA issues for a project key, resolve each to a local test function, and build slim ticket records. Populates `ticket_records` in context.
</purpose>

<prereqs>

Read [shared-jira-protocol.md](shared-jira-protocol.md) before parsing any previous-investigation comments. It defines the comment format, parsing rules, and abandonment procedure.
</prereqs>

<parent-setup>
Run `git remote get-url origin` once. Extract `{owner}/{repo}` from the URL. Cache as `current_repo`. Pass to the subagent below.
</parent-setup>

<subagent id="fetch-filter" model="haiku">

<inputs>
`KEY`, `N`, `cloudId`, `current_repo`, `nav_tool`, `lsp_available` (from `phase_outputs.phase0`), repo root path.
</inputs>

<output>
```json
{
  "slim_records": [...],
  "skipped": { "cross_repo": int, "system_tests": int, "not_found": int }
}
```
Never return raw JIRA API objects — parent only receives slim records.
</output>

<loop>
```
results = []
cursor = null
while len(results) < N:
  fetch N issues via mcp__atlassian__searchJiraIssuesUsingJql:
    jql:           project = {KEY} AND labels = "flaky-test" AND status = "Open" ORDER BY created DESC
    fields:        ["summary", "description", "comment", "status", "assignee",
                    "customfield_13010", "customfield_13009", "customfield_13007"]
    maxResults:    N
    nextPageToken: cursor  (omit on first call)

  for each issue (in order):
    1. Repo check (zero-cost): extract {owner}/{repo} from customfield_13009
       (2nd+3rd segments after github.com/). Mismatch → skip (cross_repo++).
       If customfield_13009 absent, scan description for
       https://github.com/{owner}/{repo} or a "Repo:" / "Repository:" field.

    2. System-tests exclusion (zero-cost): if customfield_13009 starts with
       github.com/smartcontractkit/chainlink/system-tests/ → skip (system_tests++).

    3. Test function check: extract top-level function name from customfield_13007
       (part before first /), fall back to longest TestXxx token in title if absent.
       - nav_tool="lsp" or lsp_available=true: LSP definition lookup
       - nav_tool="crg": mcp__code-review-graph__semantic_search_nodes_tool
       - last resort only: grep -rl "func {TestName}" .
       Not found → skip (not_found++).

    4. Eligible: build slim record (see schema below), append to results.
       Stop once len(results) == N.

  cursor = nextPageToken from response
  if no more pages: break
```
</loop>

<slim-record-schema>
```json
{
  "key": "KEY-NNN",
  "title": "...",
  "description": "...",
  "trunk_test_case_url": "https://app.trunk.io/.../test/{UUID}",
  "test_case_id": "{UUID}",
  "package": "github.com/owner/repo/...",
  "test_name": "TestFoo/subtest_name",
  "previous_attempts": [{
    "outcome": "str",
    "date": "str",
    "summary": "str",
    "excluded_approaches": ["str"],
    "rejection_reasons": ["str"],
    "recommended_next_step": "str | null",
    "full_text": "str"
  }]
}
```

Field extraction rules:
- `test_case_id`: `customfield_13010` (bare UUID). If absent, extract UUID from `https://app.trunk.io/*/test/{UUID}` in description. Null only if neither yields a value.
- `package`: `customfield_13009`. Null if absent.
- `test_name`: `customfield_13007` (full path including subtest, e.g. `TestFoo/subtest`). If absent, longest `TestXxx`/`testXxx` token from title.
- `trunk_test_case_url`: scan description for `https://app.trunk.io/*/test/{UUID}`; null if not found. Display only.
- `previous_attempts`: parse per `shared-jira-protocol.md`.
- `ci_run_url`: always `null` in project mode (no upfront URL syntax). The 3a fallback may prompt for one later.
- If any custom field is absent from the search response, call `mcp__atlassian__getJiraIssue` with `fields=["summary","description","comment","status","assignee","customfield_13010","customfield_13009","customfield_13007"]` for that issue as a fallback.
</slim-record-schema>

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
