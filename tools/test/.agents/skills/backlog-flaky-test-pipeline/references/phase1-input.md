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
4. Neither matched → **ask the user**:

<user-prompt>
**Two modes available:**

**Project mode** — searches JIRA for open flaky-test issues in a project, filters to ones that exist in this repo, and claims the top N for investigation. Good for batch-processing a queue.
- Requires: project key (e.g. `CRE`, `CCIP`, `DX`) and number of issues (default 3).

**Direct-ticket mode** — you provide specific JIRA ticket IDs (e.g. `CRE-5719 CCIP-42`) and those exact tickets are investigated. Good when you already know which tickets to fix. Optionally attach a CI run URL with `KEY@<url>` (e.g. `CRE-5719@https://github.com/.../actions/runs/123`) to feed it to `investigate-ci-failure`.

Which mode would you like? (Or just paste ticket IDs / a project key to pick implicitly.)
</user-prompt>

Once the user responds, extract `KEY` + `N` for project mode, or ticket IDs for direct-ticket mode.
</mode-detection>

<validation>
If `N > 5`: suggest a lower number and wait for confirmation before proceeding. Accept on second confirmation even if N is still high.
</validation>

<on_complete>
Write to `invocation`: `{ mode, args, auto_mode }`.

- Local mode → Read [phase2d-local-mode.md](phase2d-local-mode.md) and follow its instructions.
- Project mode → Read [phase2a-project-mode.md](phase2a-project-mode.md) and follow its instructions.
- Direct-ticket mode → Read [phase2b-direct-mode.md](phase2b-direct-mode.md) and follow its instructions.
</on_complete>

</phase>
