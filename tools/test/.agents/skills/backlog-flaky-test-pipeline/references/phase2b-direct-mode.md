---
phase: phase2b
model_tier: lightweight
---

<phase id="phase2b">

<purpose>
Validate each explicitly-provided ticket, build slim records for eligible ones, and surface ownership conflicts to the user before proceeding. Populates `ticket_records` and `user_decisions` in context.
</purpose>

<prereqs>

Read [../../_shared-jira-flaky-ops/investigation-comment.md](../../_shared-jira-flaky-ops/investigation-comment.md) before parsing any previous-investigation comments. It defines the comment format and parsing rules.
</prereqs>

<parent-setup>
Run `git remote get-url origin` once. Extract `{owner}/{repo}`. Cache as `current_repo`.
</parent-setup>

<subagent id="validation" model_tier="lightweight" instances="one per ticket" parallelism="all in a single message">

Spawn one subagent per ticket, following `_shared-jira-flaky-ops/validate-flaky-ticket.md`. Read that file for the full validation steps and field extraction rules.

Inputs to pass per ticket: ticket key, `ci_run_url` (from phase1 `KEY@URL` parsing — null if not provided), `cloudId`, `current_repo`, `nav_tool`, `lsp_available` (from `phase_outputs.phase0`), repo root path.

Expected output per subagent:
```json
{
  "status": "ok" | "error" | "needs_assignment_check",
  "message": "string (errors only)",
  "slim_record": { /* see _shared-jira-flaky-ops/README.md for schema */ },
  "assignee_display_name": "string (needs_assignment_check only)",
  "current_status": "string (needs_assignment_check only)"
}
```

</subagent>

<on_subagent_return>
Handle each result (record ownership decisions in `user_decisions`):

- `"error"` → inform the user with the message; skip this ticket.
- `"needs_assignment_check"` → surface to user and await explicit confirmation:
  > KEY-NNN is currently assigned to {displayName} and is in '{current_status}' status — someone else may already be working on it. How would you like to proceed?

  Do not claim or transition until the user explicitly confirms. On confirmation → use the `slim_record` from the result and proceed to claim. Record decision in `user_decisions`.
- `"ok"` → add `slim_record` to `ticket_records` and proceed to claim.
</on_subagent_return>

<on_complete>
Read [phase2c-prior-gate.md](phase2c-prior-gate.md) and follow its instructions.
</on_complete>

</phase>
