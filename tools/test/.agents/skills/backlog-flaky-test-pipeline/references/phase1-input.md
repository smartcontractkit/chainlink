---
phase: phase1
model_tier: lightweight
---

<phase id="phase1">

<purpose>
Resolve invocation mode and arguments. Write `invocation.mode` and `invocation.args` to context, then route to the appropriate phase 2 variant.
</purpose>

<mode-detection>
Check arguments in order — stop at the first match:

1. **`--local` flag present** → **local-test mode**. Parse remaining args as:
   - Test specs: any token not starting with `--` and not immediately following `--log`. Format: `<pkg>.TestName` (period separator) or bare `TestName`. Subtest segments use `/`.
   - `--log <path>`: the token immediately following `--log` is the log file path.
   - Validate: at least one test spec is present. If none found, re-prompt: *"Local mode requires at least one test spec. Example: `--local core/services/llo.TestFoo`"*
   - Set `args = { test_specs: ["..."], log_path: "<path> | null" }`. Skip prompt.

2. Any argument matches `PROJ-NNN` or `PROJ-NNN@<url>` (e.g. `CRE-5719`, `CCIP-42@https://github.com/.../actions/runs/123`) → **direct-ticket mode**, args = list of `{ key, ci_run_url }` pairs. For each token: split on the first `@`. Left side is the JIRA key (must match `PROJ-NNN`); right side, if present, is the CI run URL — store as-is, do not validate the URL here. Skip prompt.
3. Both `KEY` and `N` were provided → **project mode**, args = `{ key, n }`. Skip prompt.
4. Neither matched → **ask the user** using `AskUserQuestion`:
- Question: "Which mode would you like?"
- Options:
  - "Project mode — search JIRA for the top N open flaky-test tickets in a project (provide project key + count)"
  - "Direct-ticket mode — investigate specific JIRA ticket IDs you provide (e.g. CRE-5719, or CRE-5719@<ci-run-url>)"

Once the user responds, ask a follow-up for the required inputs (project key + N, or ticket IDs) if not already provided.
</mode-detection>

<validation>
If `N > 5`: suggest a lower number and wait for confirmation before proceeding. If user reconfirms, proceed without further prompting.
</validation>

<on_complete>
Write to `invocation`: `{ mode, args, auto_mode }`.

- Local mode → Read [phase2d-local-mode.md](phase2d-local-mode.md) and follow its instructions.
- Project mode → Read [phase2a-project-mode.md](phase2a-project-mode.md) and follow its instructions.
- Direct-ticket mode → Read [phase2b-direct-mode.md](phase2b-direct-mode.md) and follow its instructions.
</on_complete>

</phase>
