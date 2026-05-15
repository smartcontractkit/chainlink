---
phase: phase6
model: haiku
---

<phase id="phase6">

<purpose>
Write final JIRA comments and transition tickets to their terminal state. Reads `lint_status` and `lint_scope` from `ticket_records`.
</purpose>

<prereqs>

Read [../../_shared-jira-flaky-ops/investigation-comment.md](../../_shared-jira-flaky-ops/investigation-comment.md) if not already loaded — needed for Investigation Update comment format.

Read [../../_shared-jira-flaky-ops/abandon-ticket.md](../../_shared-jira-flaky-ops/abandon-ticket.md) if not already loaded — needed for INCONCLUSIVE abandonment procedure.
</prereqs>

<fixed-issues>
For each FIXED issue:

1. Follow `_shared-jira-flaky-ops/transition-ticket.md` with `jira_key` and `target = "In Review"`.
2. Follow `_shared-jira-flaky-ops/investigation-comment.md` to write `addCommentToJiraIssue` (OUTCOME = FIXED):
   - **What was investigated**: the failure mode and root cause in one sentence.
   - **Hypothesis**: the Proposer's root cause.
   - **What was tried**: fix description; PR: {PR URL}; classification: {classification} ({confidence}); pattern: {pattern_category}; rationale: {rationale}. If classification was SUT with user override, note it here.
   - **Why it didn't hold**: N/A.
   - **Recommended next step**: N/A.
</fixed-issues>

<inconclusive-issues>
For each INCONCLUSIVE issue:

1. Follow `_shared-jira-flaky-ops/investigation-comment.md` to write `addCommentToJiraIssue` (OUTCOME = INCONCLUSIVE):
   - **What was investigated**: the failure mode and what code was analyzed.
   - **Hypothesis**: the Proposer's root cause.
   - **What was tried**: the proposed fix if any, otherwise "No fix applied."
   - **Why it didn't hold**: the Challenger's key objections and the Arbiter's rationale.
   - **Recommended next step**: a concrete actionable direction derived from the Arbiter's reasoning.
2. Follow `_shared-jira-flaky-ops/abandon-ticket.md` (unassign + transition to "Open").
</inconclusive-issues>

<on_complete>
Print final summary:

> **Session complete.**
> - Fixed: [KEY-1, KEY-2] → PR: {PR URL}
> - Returned to queue: [KEY-3] (INCONCLUSIVE / PARTIAL_FIX)
> - No further action required unless you want to follow up on the returned tickets.
</on_complete>

</phase>
