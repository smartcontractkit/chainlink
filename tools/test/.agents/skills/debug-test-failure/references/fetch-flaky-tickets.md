<fetch-flaky-tickets>
JQL search loop that fetches N eligible flaky-test tickets for a project key, filters out cross-repo and system-test tickets, resolves each to a local test function, and returns slim records

<requirements>
1. JIRA project provided by the user.
2. Number of issues to search for.
3. Current repository — `{owner}/{repo}` extracted from `git remote get-url origin`
</requirements>

<search_loop>

```
results = []
cursor = null
while len(results) < N:
  fetch N issues via mcp__atlassian__searchJiraIssuesUsingJql:
    jql:           project = {KEY} AND labels = "flaky-test" AND status = "Open" ORDER BY created DESC
    fields:        ["summary", "description", "comment", "status", "assignee",
                    "customfield_13009", "customfield_13007"]
    maxResults:    N
    nextPageToken: cursor  (omit on first call)

  for each issue (in order):
    1. Repo check (zero-cost): extract {owner}/{repo} from customfield_13009
       (2nd + 3rd segments after github.com/). Compare with current repository (from `git remote get-url origin`).
       Mismatch → skip (cross_repo++).
       If customfield_13009 absent, scan description for
       https://github.com/{owner}/{repo} or a "Repo:" / "Repository:" field.

    2. System-tests exclusion (zero-cost): if customfield_13009 starts with
       github.com/smartcontractkit/chainlink/system-tests/ → skip (system_tests++).

    3. Test function check: extract top-level function name from customfield_13007
       (part before first /), fall back to longest TestXxx token in title if absent.
       - `LSP` available: LSP definition lookup
       - `Code Review Graph` available: mcp__code-review-graph__semantic_search_nodes_tool
       - last resort only: grep -rl "func {TestName}" .
       Not found → skip (not_found++).

    4. Eligible: build slim record (see schema), append to results.
       Stop once len(results) == N.

  cursor = nextPageToken from response
  if no more pages: break
```
</search_loop>

Claim the eligible tickets via [claim-ticket](./claim-ticket.md), skip pre-requisites (you have already checked them).

</fetch-flaky-tickets>