---
phase: phase5
model: haiku
---

<phase id="phase5">

<purpose>
Commit fixes, push, and open a PR. Updates `ticket_records` with `branch`, `commit_sha`, and `pr_url`.
</purpose>

<steps>

<step id="1-summary">
Present summary table: issue key | verdict | fix file(s) | test result.
</step>

<step id="2-branch">
Ask: commit on current branch or new? Determine the proposed branch name:
- Start with `fix/flaky-{KEY}-{YYYY-MM-DD}`.
- Check for conflicts: `git branch --list "fix/flaky-{KEY}-{YYYY-MM-DD}"` and `git branch -r --list "*/fix/flaky-{KEY}-{YYYY-MM-DD}"`.
- If taken, try `-2`, `-3`, etc. until a free name is found.
- Propose the first available name.

In `--auto` mode: use the first available name without prompting.
</step>

<step id="3-stage">
Stage only the changed files by name — never `git add .`.
</step>

<step id="4-yubikey-warning">
Before committing, inform the user:

<user-prompt>
If your git commits require hardware key signing (e.g. GPG with a YubiKey), you may need to tap your YubiKey when the commit command runs.
</user-prompt>
</step>

<step id="5-commit">
Commit following the repo's existing message style.
</step>

<step id="6-ownership-recheck" never-skip-in-auto="true">
For each FIXED issue: call `mcp__atlassian__getJiraIssue`, confirm assignee still matches cached `accountId`.

If reassigned: pause and report. Do not push until user explicitly confirms. **This gate is never skipped by `--auto`.**
</step>

<step id="7-pr-dedup">
Run `gh pr list --state open --search "{KEY}"`. If a PR already exists: skip `gh pr create`, reuse the existing PR URL for the JIRA comment.
</step>

<step id="8-push">
Push. If push requires auth: ask the user to run `! git push`.
</step>

<step id="9-pr-create">
Run `gh pr create` with a description that includes:
- One-line summary of each fix.
- **JIRA links** — direct link to each fixed ticket.
- **Trunk links** — `[Trunk analysis](trunk_analysis_url)` if available, otherwise `[Trunk test case](trunk_test_case_url)`.

Example PR body section:
```
## Flaky test fixes

| Issue | Test | Trunk |
|-------|------|-------|
| [KEY-123](jira_url) | `TestFooBar` | [Analysis](trunk_analysis_url) |
```

Capture the PR URL. Write `pr_url` and `branch` to `ticket_records`.
</step>

</steps>

<on_complete>
Announce: "PR created: {PR URL}. Moving to JIRA update."

Read [phase6-jira-update.md](phase6-jira-update.md) and follow its instructions.
</on_complete>

</phase>
