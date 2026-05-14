---
name: shared-jira-protocol
description: Investigation Update comment format, mid-flight abandonment rule, and previous-attempt parsing. Read by all phases that write JIRA comments.
---

<protocol id="shared-jira-protocol">

<comment-format>
Every `addCommentToJiraIssue` call in this skill uses the structure below. Always include all five sections — write `N/A` for any that don't apply. **Style: concise and matter-of-fact. 1–3 sentences per section. No narrative padding, no hedge words.**

```markdown
## Investigation Update — {OUTCOME} · {YYYY-MM-DD}

**Outcome**: {OUTCOME}
**Investigator**: {display name from atlassianUserInfo}
**Classification**: {TEST | SUT | AMBIGUOUS | INFRA} (confidence: {high | low | none}) | N/A

### What was investigated
{The failure mode and where analysis focused.}

### Hypothesis
{The proposed root cause, or N/A.}

### What was tried
{The fix or approach applied or proposed, or N/A.}

### Why it didn't hold
{Objections, test results, or reason the fix was rejected or reverted — or N/A.}

### Recommended next step
{Concrete actionable direction for the next investigator, or N/A.}
```

**Outcome values**: `INCONCLUSIVE` (debate unresolved), `PARTIAL_FIX` (fix applied, tests still failed, reverted), `FIXED` (fix verified, PR created), `RETURNED_TO_QUEUE` (SUT/AMBIGUOUS/INFRA classification), `CLOSED_SUBTEST` (failure in t.Run subtest, not top-level), `ABANDONED` (mid-flight stop for any reason).
</comment-format>

<abandonment-rule>
If the user cancels, skips, or stops working on a JIRA ticket at **any** point after it was claimed — regardless of reason — you **must** (never skip):

1. `mcp__atlassian__editJiraIssue` → unassign the issue (set assignee to null).
2. `mcp__atlassian__getTransitionsForJiraIssue` → find "Open".
3. `mcp__atlassian__transitionJiraIssue` → transition back to "Open".
4. `mcp__atlassian__addCommentToJiraIssue` → write an Investigation Update comment (OUTCOME = ABANDONED). "What was investigated": state the reason work stopped. All other sections: N/A.

Applies when: user says "skip this one", verdict is INCONCLUSIVE, PARTIAL_FIX is reverted, ownership conflict detected, or user ends the session early. Do **not** leave claimed tickets in "In Progress" with no assignee action.
</abandonment-rule>

<parsing-previous-attempts>
Scan JIRA comments for `## Investigation Update — {OUTCOME}`. For each match extract:

- `outcome`: the OUTCOME token from the heading
- `date`: the date after `·`
- `full_text`: the full comment text
- `excluded_approaches`: content of `### What was tried` (skip if "N/A")
- `rejection_reasons`: content of `### Why it didn't hold` (skip if "N/A")
- `recommended_next_step`: content of `### Recommended next step` (null if "N/A")
- `summary`: content of `### What was investigated`

Fall back to keyword scanning for non-standard-format comments.
</parsing-previous-attempts>

</protocol>
