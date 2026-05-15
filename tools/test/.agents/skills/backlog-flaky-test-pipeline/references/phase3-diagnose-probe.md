---
name: phase3-diagnose-probe
description: Shared procedure — invoke the chainlink diagnose tool to gather local failure symptoms. Called when no log or Trunk data is available. Writes to actionable_facts and trunk_investigation_status.
---

<procedure id="diagnose-probe">

<purpose>
Run the chainlink `diagnose` tool once to collect failure symptoms locally. Evidence-gathering mode only. Writes to `actionable_facts` and `trunk_investigation_status` in the calling context.
</purpose>

<inputs>
- `test_name` — top-level test function name (part before the first `/` in `slim_record.test_name`).
- `package` — Go package path (e.g. `core/services/llo`).
- `caller_context` — label for user-facing messages (JIRA key or test name).
</inputs>

<invocation>
```bash
go -C tools/test run . diagnose --ai-output --iterations 10 --parallel-iterations 5 -- --run "^{test_name}$" --race --shuffle=on ./{package}/...
```

Drop `--parallel-iterations` to `1` if the test holds external resources (fixed port, exclusive temp file, exclusive lock) that cannot tolerate concurrent runs.

Package-scope (`./{package}/...`) is intentional even though `-run` limits which tests execute: Go still loads the full package, so `TestMain` / `init()` / package-level setup runs, surfacing init-order and shared-state flakes. Cross-test ordering effects (TestY pollutes state, TestX reads it) will not surface because only the named test executes — note this to the user if all iterations pass.
</invocation>

<outcome-handling>
Parse the `--ai-output` summary:

- **At least one iteration failed** → extract failure-specific portions (error messages, stack traces, race-detector output, timeout reports) into `actionable_facts` as raw strings. Set `trunk_investigation_status = "diagnose_run"`.
- **All iterations passed** → `actionable_facts` stays `[]`, `trunk_investigation_status` unchanged. Inform user: *"`diagnose` ran 10 iterations without reproducing the failure for {caller_context} (single-test scope — cross-test ordering effects won't surface); proceeding with code analysis only."*
- **Tool failed to run** (missing dependency, build error, etc.) → `actionable_facts` stays `[]`, status unchanged. Log the error; do not retry.
</outcome-handling>

</procedure>
