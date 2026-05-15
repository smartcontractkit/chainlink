---
name: shared-jira-flaky-ops-readme
description: Index of shared JIRA operations for flaky-test skills. Defines the canonical slim-record schema, required env, and how to use these files from another skill.
---

# Shared JIRA Flaky-Test Operations

This directory contains self-contained reference files for JIRA operations specific to flaky-test workflows. Any skill that needs to interact with JIRA around flaky tests should read the relevant file here instead of duplicating the logic.

## When to use

These files are **includable references**, not user-facing skills (no `SKILL.md`). Read the relevant file(s) when you need to perform one of the operations listed below. Each file is self-contained: it declares its inputs, steps, and outputs. You do not need to read files you are not using.

## Required environment

All operations require:
- `cloudId` — the Atlassian cloud ID (from `mcp__atlassian__getAccessibleAtlassianResources` or cached from a prior call).
- `accountId` — the current user's Atlassian account ID (from `mcp__atlassian__atlassianUserInfo`, cached in `phase_outputs.phase0`).

## Operations index

| File | Purpose | Key inputs |
|------|---------|------------|
| `investigation-comment.md` | Comment format for Investigation Updates; parsing prior-attempt comments | jira_key, outcome, field values |
| `abandon-ticket.md` | Mid-flight abandonment: unassign → Open → Investigation Update comment (ABANDONED) | jira_key, reason |
| `transition-ticket.md` | Transition a ticket to a semantic target state | jira_key, target |
| `claim-ticket.md` | Assign to self and transition to "In Progress" | jira_key, accountId |
| `fetch-flaky-tickets.md` | JQL search loop: fetch N eligible flaky-test tickets for a project key | KEY, N, cloudId, current_repo, nav_tool, lsp_available, repo root |
| `validate-flaky-ticket.md` | Validate a single explicitly-provided ticket and build its slim record | jira_key, ci_run_url, cloudId, current_repo, nav_tool, lsp_available, repo root |
| `recheck-ownership.md` | Verify the ticket is still assigned to us before touching files or pushing | jira_key, accountId |

## How to include from another skill

1. Read the relevant file(s) from this directory.
2. Collect the declared inputs from your skill's context.
3. Follow the file's steps exactly.
4. Use the declared output schema to integrate the result back into your skill's state.

Example: to claim a ticket, read `claim-ticket.md` and follow its steps with your `jira_key` and `accountId`.

---

## Canonical slim-record schema

Defined here once. Referenced by `fetch-flaky-tickets.md` and `validate-flaky-ticket.md`. Both files must produce records conforming to this schema. Do not redefine it in either file.

```json
{
  "jira_key":            "KEY-NNN | null",
  "local_id":            "local-N | null",
  "title":               "string",
  "description":         "string",
  "trunk_test_case_url": "https://app.trunk.io/.../test/{UUID} | null",
  "test_case_id":        "{UUID} | null",
  "package":             "github.com/owner/repo/path | null",
  "test_name":           "TestFoo | TestFoo/subtest_name",
  "previous_attempts":   [
    {
      "outcome":               "INCONCLUSIVE | PARTIAL_FIX | MISMATCH | SKIP_TOP_LEVEL | RETURNED_TO_QUEUE | ABANDONED | FIXED",
      "date":                  "YYYY-MM-DD",
      "summary":               "string",
      "excluded_approaches":   ["string"],
      "rejection_reasons":     ["string"],
      "recommended_next_step": "string | null",
      "full_text":             "string"
    }
  ],
  "ci_run_url":          "string | null",
  "provided_log_path":   "string | null",
  "provided_log_text":   "string | null"
}
```

**Field rules:**
- `jira_key`: non-null in JIRA modes; null in local mode.
- `local_id`: non-null in local mode (`local-1`, `local-2`, …); null in JIRA modes.
- `test_case_id`: `customfield_13010` (bare UUID). If absent, extract UUID from `https://app.trunk.io/*/test/{UUID}` in description. Null only if neither yields a value. Always null in local mode.
- `package`: `customfield_13009`. Null if absent.
- `test_name`: `customfield_13007` (full path including subtest). If absent, longest `TestXxx`/`testXxx` token from title.
- `trunk_test_case_url`: scan description for `https://app.trunk.io/*/test/{UUID}`; null if not found. Display only.
- `previous_attempts`: parsed per `investigation-comment.md` parsing rules. Empty array in local mode.
- `ci_run_url`: null in project mode (no upfront URL syntax). Populated from `KEY@URL` syntax in direct-ticket mode, or from a 3a fallback prompt.
- `provided_log_path`, `provided_log_text`: null in JIRA modes; populated from `--log <path>` in local mode.
