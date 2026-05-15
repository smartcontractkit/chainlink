---
name: validate-flaky-ticket
description: Validates a single explicitly-provided JIRA ticket key and builds its slim record. Returns structured status so the caller can handle ownership conflicts. Extracted from phase2b.
---

# Validate Flaky Ticket

## Inputs

- `jira_key` — the JIRA ticket key (e.g. `CRE-5719`)
- `ci_run_url` — GitHub Actions run URL from `KEY@URL` syntax, or null if not provided
- `cloudId` — Atlassian cloud ID
- `current_repo` — `{owner}/{repo}` extracted from `git remote get-url origin`
- `nav_tool` — `"lsp"` | `"crg"` (from `phase_outputs.phase0`)
- `lsp_available` — boolean (from `phase_outputs.phase0`)
- repo root path

## Output

```json
{
  "status": "ok" | "error" | "needs_assignment_check",
  "message": "string (errors only)",
  "slim_record": { /* see README.md slim-record schema — ok and needs_assignment_check only */ },
  "assignee_display_name": "string (needs_assignment_check only)",
  "current_status": "string (needs_assignment_check only)"
}
```

Never return raw JIRA API objects. The slim record schema is defined in `README.md`.

## Steps

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

6. **Build slim record** from the step 1 response — no additional API call.
   - `jira_key` non-null, `local_id` null.
   - `ci_run_url` from input parameter (null if not provided).
   - `provided_log_path`, `provided_log_text` null.
   - `previous_attempts`: parse per `investigation-comment.md` parsing rules.
   - Return `{ "status": "ok", "slim_record": {...} }`.
