---
phase: phase6
model: haiku
---

<phase id="phase6">

<purpose>
Write final JIRA comments and transition tickets to their terminal state. Reads `lint_status` and `lint_scope` from `ticket_records`.
</purpose>

<prereqs>

Read [shared-jira-protocol.md](shared-jira-protocol.md) if not already loaded in this session.
</prereqs>

<fixed-issues>
For each FIXED issue:

1. `getTransitionsForJiraIssue` → find "In Review" (aliases: "In Code Review", "Review").
2. `transitionJiraIssue` → "In Review".
3. `addCommentToJiraIssue` → Investigation Update comment (OUTCOME = FIXED):
   - **What was investigated**: the failure mode and root cause in one sentence.
   - **Hypothesis**: the Proposer's root cause.
   - **What was tried**: fix description; PR: {PR URL}; signals matched: {sut_signals_matched + test_signals_matched}; SUT score: {sut_score}, TEST score: {test_score}. If classification was SUT with user override, note it here.
   - **Why it didn't hold**: N/A.
   - **Recommended next step**: N/A.
</fixed-issues>

<inconclusive-issues>
For each INCONCLUSIVE issue:

1. `addCommentToJiraIssue` → Investigation Update comment (OUTCOME = INCONCLUSIVE):
   - **What was investigated**: the failure mode and what code was analyzed.
   - **Hypothesis**: the Proposer's root cause.
   - **What was tried**: the proposed fix if any, otherwise "No fix applied."
   - **Why it didn't hold**: the Challenger's key objections and the Arbiter's rationale.
   - **Recommended next step**: a concrete actionable direction derived from the Arbiter's reasoning.
2. Apply the mid-flight abandonment rule (unassign + transition to "Open").
</inconclusive-issues>

<on_complete>
Print final summary:

> **Session complete.**
> - Fixed: [KEY-1, KEY-2] → PR: {PR URL}
> - Returned to queue: [KEY-3] (INCONCLUSIVE / PARTIAL_FIX)
> - No further action required unless you want to follow up on the returned tickets.
</on_complete>

</phase>
