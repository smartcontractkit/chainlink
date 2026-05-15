---
name: claim-ticket
description: Assign a flaky-test ticket to the current user and transition it to In Progress. Used after a ticket is approved for investigation.
---

# Claim Ticket

## Inputs

- `jira_key` — the JIRA ticket key (e.g. `CRE-5719`)
- `accountId` — the current user's Atlassian account ID (from `phase_outputs.phase0`)

## Steps

Execute in order — wait for each step to succeed before proceeding:

1. `mcp__atlassian__getJiraIssue` with `jira_key` → read `fields.assignee.accountId`. Save it as `original_assignee` (null if the field is absent or the ticket is unassigned).
2. `mcp__atlassian__editJiraIssue` → assign the issue to `accountId` (set `assignee.accountId = accountId`). Wait for success.
3. Follow `transition-ticket.md` with `jira_key` and `target = "In Progress"`.
   - If the transition fails: log available transitions and stop. Do not leave the ticket assigned without transitioning.

## Output

```json
{ "success": true, "jira_key": "KEY-NNN", "original_assignee": "<accountId or null>" }
```

or on failure:

```json
{ "success": false, "jira_key": "KEY-NNN", "error": "<reason>" }
```
