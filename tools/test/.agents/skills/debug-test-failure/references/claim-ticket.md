<claim-ticket>
Assign a flaky-test ticket to the current user and transition it to In Progress.

<pre-requisites>
1. Ticket status check: claim if `status == "Open"` OR `(status != "Open" AND assignee == current user)`.
    Otherwise, do not claim the ticket; in `skip_reason` say that someone else is already working on it.

2. Repo check (zero-cost): extract {owner}/{repo} from `package` field
    (2nd + 3rd segments after github.com/). Compare with current repository (extracted from `git remote get-url origin`)
    Mismatch → do not claim the ticket, in `skip_reason` say that this ticket should be worked on
    in the context of the {repo}.

3. System-tests exclusion (zero-cost): if `package` field starts with
    github.com/smartcontractkit/chainlink/system-tests/ → do not claim the ticket,
    in `skip_reason` say that system-tests are excluded from the skill due to their
    complexity and token cost.

4. Test function check: extract top-level function name from `test_name` field
    (part before first /), fall back to longest TestXxx token in title if absent.
    - `LSP` available: LSP definition lookup
    - `Code Review Graph` available: `mcp__code-review-graph__semantic_search_nodes_tool`
    - last resort only: grep -rl "func {TestName}" .
    Not found → do not claim the ticket, in `skip_reason` say that this test doesn't exist in {repo}.
</pre-requisites>

<steps>
Execute in order — wait for each step to succeed before proceeding:

1. `mcp__atlassian__getJiraIssue` with `jira_key` → read `fields.assignee.accountId`. Save it as `original_assignee` (null if the field is absent or the ticket is unassigned).
2. `mcp__atlassian__editJiraIssue` → assign the issue to `accountId` (set `assignee.accountId = accountId`). Wait for success.
3. Unless the ticket is already assigned to current user and is in `In Progress` stage follow [transition-ticket.md](./transition-ticket.md) with `jira_key` and `target = "In Progress"`.
   - If the transition fails: unassign the ticket.
   - Set transition failure in: `skip_reason`.
</steps>

</claim-ticket>