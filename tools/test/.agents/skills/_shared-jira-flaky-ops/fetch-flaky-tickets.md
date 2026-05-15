---
name: fetch-flaky-tickets
description: JQL search loop that fetches N eligible flaky-test tickets for a project key, filters out cross-repo and system-test tickets, resolves each to a local test function, and returns slim records. Extracted from phase2a.
---

# Fetch Flaky Tickets

## Inputs

- `KEY` — JIRA project key (e.g. `CRE`)
- `N` — number of eligible records to return
- `cloudId` — Atlassian cloud ID
- `current_repo` — `{owner}/{repo}` extracted from `git remote get-url origin`
- `nav_tool` — `"lsp"` | `"crg"` (from `phase_outputs.phase0`)
- `lsp_available` — boolean (from `phase_outputs.phase0`)
- repo root path

## Output

```json
{
  "slim_records": [ /* see README.md slim-record schema */ ],
  "skipped": { "cross_repo": 0, "system_tests": 0, "not_found": 0 }
}
```

Never return raw JIRA API objects — caller only receives slim records.

## Slim-record schema

See `README.md` in this directory for the canonical schema. In the records produced here:
- `jira_key` non-null, `local_id` null.
- `ci_run_url` always null (project mode has no upfront URL syntax; 3a fallback may prompt later).
- `provided_log_path`, `provided_log_text` null.

## Loop

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
       (2nd + 3rd segments after github.com/). Mismatch → skip (cross_repo++).
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

    4. Eligible: build slim record (see schema), append to results.
       Stop once len(results) == N.

  cursor = nextPageToken from response
  if no more pages: break
```

## Field extraction rules

- `test_case_id`: `customfield_13010` (bare UUID). If absent, extract UUID from `https://app.trunk.io/*/test/{UUID}` in description. Null only if neither yields a value.
- `package`: `customfield_13009`. Null if absent.
- `test_name`: `customfield_13007` (full path including subtest, e.g. `TestFoo/subtest`). If absent, longest `TestXxx`/`testXxx` token from title.
- `trunk_test_case_url`: scan description for `https://app.trunk.io/*/test/{UUID}`; null if not found. Display only.
- `previous_attempts`: parse per `investigation-comment.md` parsing rules.
- If any custom field is absent from the search response, call `mcp__atlassian__getJiraIssue` with `fields=["summary","description","comment","status","assignee","customfield_13010","customfield_13009","customfield_13007"]` for that issue as a fallback.
