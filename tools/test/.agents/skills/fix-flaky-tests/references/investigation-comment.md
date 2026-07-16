<investigation-comment>
<writing_a_comment>
Every `mcp__atlassian__addCommentToJiraIssue` call for an investigation update uses the structure below. Always include all five sections — write `N/A` for any that don't apply. **Style: concise and matter-of-fact. 1–3 sentences per section. No narrative padding, no hedge words.**

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
</writing_a_comment>

<parsing_previous-attempt_comments>
Scan JIRA comments for `## Investigation Update — {OUTCOME}`. For each match extract:

- `outcome` — the OUTCOME token from the heading
- `date` — the date after `·`
- `full_text` — the full comment text
- `excluded_approaches` — content of `### What was tried` (skip if "N/A")
- `rejection_reasons` — content of `### Why it didn't hold` (skip if "N/A")
- `recommended_next_step` — content of `### Recommended next step` (null if "N/A")
- `summary` — content of `### What was investigated`

Fall back to keyword scanning for non-standard-format comments.
</parsing_previous-attempt_comments>
</investigation-comment>