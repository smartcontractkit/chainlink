---
phase: phase3a-iii
model_tier: lightweight
---

<substep id="3a-iii">

<purpose>
Detect tickets where the failure originates in a `t.Run` subtest rather than the top-level test function. These should be filed against the specific subtest, not the outer function.
</purpose>

<check>
Inspect `actionable_facts` (and any stack trace within them) for a `file:line`. If the failure line falls inside a `t.Run(...)` callback AND the outer function contains no assertions outside `t.Run` blocks → candidate for **SKIP_TOP_LEVEL**.

**Exception**: if `slim_record.test_name` contains `/` the ticket was already filed against a specific subtest — SKIP_TOP_LEVEL must not fire. If `test_name` is null, check the title for `/`.
</check>

<if-skip-top-level>

**JIRA mode** (mode ≠ "local"):
1. Follow `_shared-jira-flaky-ops/transition-ticket.md` with `jira_key` and `target = "Won't Do"`.
2. Follow `_shared-jira-flaky-ops/investigation-comment.md` — write `addCommentToJiraIssue` (OUTCOME = CLOSED_SUBTEST). "What was investigated": failure originates in a `t.Run` subtest, not the top-level function. "Recommended next step": file or locate a ticket for the specific subtest. All other sections: N/A.
3. Return verdict `SKIP_TOP_LEVEL` for this issue.

**Local mode** (mode = "local"):
- Print: *"Skipping {test_name} — failure originates in a t.Run subtest, not the top-level function. File a ticket against the specific subtest if you want this fixed."*
- Return verdict `SKIP_TOP_LEVEL`. No JIRA writes.

</if-skip-top-level>

<on_complete>
If SKIP_TOP_LEVEL was not returned: proceed to `<substep id="3b">` in [phase3-investigation.md](phase3-investigation.md).
</on_complete>

</substep>
