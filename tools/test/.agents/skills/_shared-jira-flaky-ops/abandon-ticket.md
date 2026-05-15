---
name: abandon-ticket
description: Mid-flight abandonment procedure for a claimed flaky-test ticket. Unassigns, transitions back to Open, and writes an ABANDONED Investigation Update comment. Must never be skipped once a ticket has been claimed.
---

# Abandon Ticket

Apply whenever a claimed ticket is stopped mid-flight — regardless of reason. This includes: user cancels, user skips, verdict is INCONCLUSIVE, PARTIAL_FIX is reverted, ownership conflict detected, SUT/AMBIGUOUS/INFRA auto-queue return, or session ends early.

**Never leave a claimed ticket in "In Progress" with no assignee action.**

## Inputs

- `jira_key` — the JIRA ticket key
- `reason` — one sentence describing why work stopped (used in "What was investigated" section)
- `accountId` — the current user's Atlassian account ID (to confirm we own it before unassigning)

## Steps

Execute in order:

1. `mcp__atlassian__editJiraIssue` → unassign the issue (set `assignee` to null).
2. Follow `transition-ticket.md` with `jira_key` and `target = "Open"`.
3. Follow `investigation-comment.md` to write an `addCommentToJiraIssue` call:
   - **Outcome**: ABANDONED
   - **What was investigated**: `reason` (the reason work stopped).
   - **Hypothesis**: N/A
   - **What was tried**: N/A
   - **Why it didn't hold**: N/A
   - **Recommended next step**: N/A

## Output

No structured output. Caller continues with other issues after this completes.
