---
name: recheck-ownership
description: Verify that a claimed ticket is still assigned to the current user before touching files or pushing. Used in phase4 (before applying fixes) and phase5 (before pushing).
---

# Recheck Ownership

## Inputs

- `jira_key` — the JIRA ticket key
- `accountId` — the current user's Atlassian account ID (cached from `phase_outputs.phase0`)

## Steps

1. Call `mcp__atlassian__getJiraIssue` with `fields=["assignee"]`.
2. Compare `assignee.accountId` to `accountId`.

## Output

```json
{ "result": "ok" }
```

or if reassigned:

```json
{ "result": "reassigned", "reassigned_to": "<displayName>" }
```

## Caller responsibility

If `result = "reassigned"`:
- Report: *"KEY-NNN is now assigned to {displayName} — reach out before proceeding."*
- Follow `abandon-ticket.md` for this ticket.
- Continue with remaining issues.
