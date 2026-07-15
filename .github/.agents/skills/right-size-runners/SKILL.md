---
name: right-size-runners
description: Properly size CI runners for price, performance, and stability.
disable-model-invocation: true
---

<initialization>
Missing info? Ask user:
1. Workflow to optimize.
2. Priorities (default: stability > speed > cost).
3. "Speed" definition (e.g., cache hits vs raw execution).
4. "Stability" definition (e.g., are flaky tests acceptable?).
5. Determine "spot" setting. Don't ask user directly about what spot setting to use, ask questions to lead you to correct setting. [runs-on spot reference](https://runs-on.com/docs/costs/spot-pricing/)

Setup:
1. Read workflow. Identify jobs, runners, critical path, and bottlenecks.
2. Ask to optimize specific job or whole workflow.
3. Modify workflow for testing (bypass gates, use mock inputs, add `workflow_dispatch`).
4. Init or resume trial log at `.github/.agents/skills/right-size-runners/trials/<workflow>.md`.
5. Run a baseline trial with the current runner configuration to establish a performance and stability benchmark.
</initialization>

<constraints>
- OOM/Out of Disk Space failures are NEVER acceptable.
- Prefer default `ubuntu-latest`. Use `runs-on` for more resources.
- Validate runner config via available runners API before use.
- Use `gh` CLI for workflow execution and PRs.
- Only change one variable per trial to accurately assess its impact.
</constraints>

<resources>
- [runs-on docs](https://runs-on.com/docs/)
- [available runners](https://go.runs-on.com/api)
- [GitHub Action docs](https://docs.github.com/en/actions)
</resources>

<loop>
1. Define trials. Update `<workflow>.md` log.
2. Ask user to approve trials. Check if they want to run them in parallel or sequentially.
3. New (disposable) branch + commit + push. Message: `cpu=X/ram=Y`. PR title: `[DO NOT MERGE] Trial: <workflow-name>` description: details of trial
4. Trigger workflow. Wait in background without polluting context:
   ```sh
   gh run watch <run_id> -i 60 --exit-status > /dev/null 2>&1
   ```
5. Once background task completes, check final status and get job IDs:
   ```sh
   gh run view <run_id> --json conclusion,url,jobs
   ```
6. Extract runner telemetry (CPU/RAM). DO NOT fetch full logs! Extract only the "Complete runner" step:
   ```sh
   gh run view <run_id> --job <job_id> --log | awk -F'\t' '$2=="Complete runner"{print substr($0, index($0, $3))}'
   ```
7. Analyze results. Update trial log. Present summary and propose next steps.

<trial-template>
| Runner | Experiment | Expectation | Branch | Run ID | Commit | Stability | Runtime | Cost | Notes |
|---|---|---|---|---|---|---|---|---|---|
| `config` | What's tested | Expected result | `branch` | `run-id` | `sha` | Pass/Fail/Flaky | mm:ss | $ | Findings |
</trial-template>
</loop>

<complete>
When the user says stop, or all possible experiments have been exhausted:

* Suggest final recommendations.
* Summarize all findings (cost + speed + stability changes) per workflow/job as table(s) for PR description.

## [Workflow/Job Name] Runner Changes
| Approach | Runner | Stability | Runtime | Runtime Delta (Abs/%) | Cost | Cost Delta (Abs/%) |
|---|---|---|---|---|---|---|
| Old | [original runner before trials] | Pass/Fail/Flaky | mm:ss | +0:00 (+0%) | $ | +$ (+0%) |
| New | [new runner] | Pass/Fail/Flaky | mm:ss | +0:00 (+0%) | $ | +$ (+0%) |

* Remove all debugs, and make final edits on a new branch after approval, and ask user to commit.
* Cleanup all trial branches, logs, and PRs.
</complete>