<abandon-ticket>
Apply whenever a claimed ticket is stopped mid-flight — regardless of reason. This includes: user cancels, user skips, no clear or working fix found, ownership conflict detected or session ends early.

<requirements>
**Never leave a claimed ticket in "In Progress".**
</requirements>

<steps>
1. `mcp__atlassian__editJiraIssue` → unassign the issue (set `assignee` to null).
2. Follow [transition-ticket](./transition-ticket.md) with `target = "Open"`.
3. Follow [investigation-comment](./investigation-comment.md) to write an `addCommentToJiraIssue` call:
   - **Outcome**: ABANDONED
   - **What was investigated**: `reason` (the reason work stopped).
   - **Hypothesis**: N/A
   - **What was tried**: N/A
   - **Why it didn't hold**: N/A
   - **Recommended next step**: N/A
</steps>

</abandon-ticket>