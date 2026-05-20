<transition-ticket>
Generic transition operation for flaky-test JIRA tickets. Takes a semantic target state, resolves it to the actual transition name, and applies it.

<steps>
1. Call `mcp__atlassian__getTransitionsForJiraIssue` with `jira_key`.
2. Match `target` to an available transition using the alias table below. Pick the **first alias that appears** in the response.
3. Call `mcp__atlassian__transitionJiraIssue` with the matched transition ID.
   - For closing targets (`"Won't Do"`, `"Done"`): if the transition supports a `resolution` field, set `resolution = "Won't Do"` (fallback: `"Won't Fix"`) or
4. **If `target = "In Review"` and `original_assignee` and `accountId` are provided:**
   - If `original_assignee` is `null` or equals `accountId` → call `mcp__atlassian__editJiraIssue` to unassign the ticket (set `assignee` to `null`).
   - If `original_assignee` differs from `accountId` → call `mcp__atlassian__editJiraIssue` to assign the ticket back to `original_assignee` (set `assignee.accountId = original_assignee`).
</steps>

<target_alias_table>
| Semantic target | Try these names in order |
|----------------|--------------------------|
| `In Progress` | "In Progress", "In Development", "Active", "Start Progress" |
| `In Review` | "In Review", "In Code Review", "Code Review", "Review" |
| `Open` | "Open", "Reopen", "Backlog", "To Do", "Reopened" |
| `Won't Do` | "Won't Do", "Won't Fix", "Reject" |
| `Done` | "Done", "Closed", "Resolved", "Close", "Resolve" |

If no alias matches the available transitions, return an error rather than silently picking an unrelated state.
</target_alias_table>

</transition-ticket>