---
name: transition-ticket
description: Generic transition operation for flaky-test JIRA tickets. Takes a semantic target state, resolves it to the actual transition name, and applies it. Referenced by claim-ticket.md and abandon-ticket.md.
---

# Transition Ticket

## Inputs

- `jira_key` — the JIRA ticket key (e.g. `CRE-5719`)
- `target` — semantic state: `"In Progress"` | `"In Review"` | `"Open"` | `"Won't Do"` | `"Done"`
- `original_assignee` _(optional)_ — the `accountId` that was assigned before `claim-ticket.md` ran, or `null` if the ticket was unassigned. Only used when `target = "In Review"`.
- `accountId` _(optional)_ — the current user's Atlassian account ID. Only used when `target = "In Review"`.

## Steps

1. Call `mcp__atlassian__getTransitionsForJiraIssue` with `jira_key`.
2. Match `target` to an available transition using the alias table below. Pick the **first alias that appears** in the response.
3. Call `mcp__atlassian__transitionJiraIssue` with the matched transition ID.
   - For closing targets (`"Won't Do"`, `"Done"`): if the transition supports a `resolution` field, set `resolution = "Won't Do"` (fallback: `"Won't Fix"`).
4. **If `target = "In Review"` and `original_assignee` and `accountId` are provided:**
   - If `original_assignee` is `null` or equals `accountId` → call `mcp__atlassian__editJiraIssue` to unassign the ticket (set `assignee` to `null`).
   - If `original_assignee` differs from `accountId` → call `mcp__atlassian__editJiraIssue` to assign the ticket back to `original_assignee` (set `assignee.accountId = original_assignee`).
5. Output: `{ "success": true, "transition_name_used": "<actual name>" }` on success, or `{ "success": false, "error": "<reason>", "available_transitions": ["..."] }` if no alias matched.

## Target alias table

| Semantic target | Try these names in order |
|----------------|--------------------------|
| `In Progress` | "In Progress", "In Development", "Active", "Start Progress" |
| `In Review` | "In Review", "In Code Review", "Code Review", "Review" |
| `Open` | "Open", "Reopen", "Backlog", "To Do", "Reopened" |
| `Won't Do` | "Won't Do", "Won't Fix", "Reject", "Close", "Done" |
| `Done` | "Done", "Closed", "Resolved", "Close", "Resolve" |

## Error handling

If no alias matches any available transition: log all available transition names and return `success: false`. The caller decides how to proceed (stop, skip, or pick manually).
