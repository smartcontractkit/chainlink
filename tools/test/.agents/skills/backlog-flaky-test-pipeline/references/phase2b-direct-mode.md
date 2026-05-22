---
phase: phase2b
model: haiku
---

<phase id="phase2b">

<purpose>
Validate each explicitly-provided ticket, build slim records for eligible ones, and surface ownership conflicts to the user before proceeding. Populates `ticket_records` and `user_decisions` in context.
</purpose>

<prereqs>

Read [shared-jira-protocol.md](shared-jira-protocol.md) before parsing any previous-investigation comments. It defines the comment format, parsing rules, and abandonment procedure.
</prereqs>

<parent-setup>
Run `git remote get-url origin` once. Extract `{owner}/{repo}`. Cache as `current_repo`.
</parent-setup>

<subagent id="validation" model="haiku" instances="one per ticket" parallelism="all in a single message">

<inputs>
Ticket key, optional `ci_run_url` (from phase1 `KEY@URL` parsing — null if not provided), `cloudId`, `current_repo`, `nav_tool`, `lsp_available` (from `phase_outputs.phase0`), repo root path.
</inputs>

<output>
```json
{
  "status": "ok" | "error" | "needs_assignment_check",
  "message": "string (errors only)",
  "slim_record": { ... },
  "assignee_display_name": "string (needs_assignment_check only)",
  "current_status": "string (needs_assignment_check only)"
}
```
Never return raw JIRA API objects. The slim record schema is identical to phase2a.
</output>

<steps>
1. **Existence check**: Call `mcp__atlassian__getJiraIssue` with `fields=["summary","description","comment","status","assignee","customfield_13010","customfield_13009","customfield_13007"]` (native string array). Reuse this single response for all subsequent steps.
   - Issue not found → `{ "status": "error", "message": "Issue KEY-NNN not found — the project may not exist or the ticket number is invalid." }`

2. **Required data check**: The issue must contain both:
   - A test function name: `customfield_13007` non-null, OR a `TestXxx`/`testXxx` token in title or description.
   - A `testCaseId`: `customfield_13010` non-null, OR a Trunk URL matching `https://app.trunk.io/*/test/{UUID}` in description.
   - Either missing → `{ "status": "error", "message": "KEY-NNN is missing required data: {test name | Trunk ID | both}. Cannot reliably investigate without it." }`

3. **Repo compatibility check** (stop at first definitive result):
   a. Read `customfield_13009`. If present, extract `{owner}/{repo}` from 2nd + 3rd path segments after `github.com/`. Mismatch → `{ "status": "error", "message": "KEY-NNN specifies repo '{owner}/{repo}' which does not match the current repository." }`. Only fall back to scanning description for `https://github.com/{owner}/{repo}` or a `Repo:`/`Repository:` field if `customfield_13009` is absent.
   b. **Test function check**: use top-level function from `customfield_13007` (part before first `/`), falling back to longest `TestXxx`/`testXxx` token in title. Check locally: LSP → code-review-graph → grep (last resort only).
      - Not found → `{ "status": "error", "message": "Test function not found in the current repository for KEY-NNN." }`

4. **System-tests exclusion**: If `customfield_13009` starts with `github.com/smartcontractkit/chainlink/system-tests/` → `{ "status": "error", "message": "KEY-NNN is in the system-tests package ({package}), which is excluded from automated investigation." }`

5. **Assignment check**: If the issue is assigned to another user AND status ≠ `Open` → return `{ "status": "needs_assignment_check", "assignee_display_name": "...", "current_status": "...", "slim_record": {...} }`. Include the slim record so the parent can proceed without re-fetching if the user confirms.

6. **Build slim record** from the step 1 response — no additional API call. Field extraction rules are identical to phase2a. Set `slim_record.ci_run_url` from the input parameter (null if not provided). Return `{ "status": "ok", "slim_record": {...} }`.
</steps>

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
