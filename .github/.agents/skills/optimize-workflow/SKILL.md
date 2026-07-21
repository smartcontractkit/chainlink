---
name: optimize-workflow
description: Optimize CI runners and workflow structure for speed, simplicity, stability, and cost.
disable-model-invocation: true
---

<initialization>
Prompt user if missing:
1. Workflow to optimize.
2. Priorities (default: stability > speed > cost).
3. "Speed" definition (e.g., cache hits vs raw execution).
4. "Stability" definition (e.g., are flaky tests acceptable?).
5. Spot setting (infer via scenarios, do not ask directly). [runs-on spot reference](https://runs-on.com/docs/costs/spot-pricing/)
6. Acceptable structural changes (e.g., job splitting, caching, lint reorganization, trigger optimizations).

Setup:
1. Analyze workflow: identify jobs, runners, step dependencies, critical path, caching, and bottlenecks.
2. Identify structural optimizations, including but not limited to:
   - Parallelize sequential steps/jobs (e.g., matrix builds).
   - Add/fix caching (e.g., go modules, cargo, npm, docker layers).
   - Move checks/lints before compilation/tests to fail early.
   - Refactor triggers/scope.
3. Permitted scope: High freedom. OK to introduce breaking edits to internal `smartcontractkit` actions/workflows to improve speed/simplicity.
4. Ask user to target: specific job, whole workflow runner config, or workflow structure.
5. Setup test workflow (bypass gates, mock inputs, add `workflow_dispatch`).
6. Init/resume trials log: `.github/.agents/skills/optimize-workflow/trials/[workflow]/summary.md`.
7. Run/find/ask for a baseline trial with current configuration to benchmark against.
</initialization>

<constraints>
- No OOM or Out of Disk Space failures allowed.
- Prefer default `ubuntu-latest`. Use `runs-on` only for larger runner requirements.
  - Note: `runs-on` with `tmpfs` enabled can cause silent OOMs (not visible in memory metrics) under high disk load. Toggle `tmpfs` off if unexplained OOMs occur, and evaluate speed/cost impact.
- Verify runner configuration via available runners API before use.
- Validate workflow YAML syntax (e.g., `actionlint` or dry-run) before committing.
- Maintain semantic correctness. OK to break internal `smartcontractkit` action/workflow interfaces for speed/simplicity.
- Use `gh` CLI for runs and PRs.
- Change only one variable per trial.
- Ensure fair comparison: cache status must be identical across trials compared for speed.
- Document runner config and structure layout in each trial to ensure reproducibility.
</constraints>

<resources>
- [runs-on docs](https://runs-on.com/docs/)
- [available runners](https://go.runs-on.com/api)
- [GitHub Action docs](https://docs.github.com/en/actions)
</resources>

<loop>
1. Define trials. Update `<workflow>.md` log.
2. User approves trials and execution method (parallel vs sequential).
3. Commit + push to new disposable branch, `trial/[trial-name]`.
4. Trigger workflow.
5. Monitor run:
   `python3 .github/.agents/skills/optimize-workflow/scripts/workflow_monitor.py [run_id] [trial-name]`
   - Creates `.github/.agents/skills/optimize-workflow/trials/[workflow]/[trial-name]/report.json`
   - Creates `.github/.agents/skills/optimize-workflow/trials/[workflow]/[trial-name]/report.md`
   - Extracts runner logs to `.github/.agents/skills/optimize-workflow/trials/[workflow]/[trial-name]/logs/`
6. Compare trials:
   `python3 .github/.agents/skills/optimize-workflow/scripts/workflow_compare.py .github/.agents/skills/optimize-workflow/trials/[workflow]/[trial-1]/report.json .github/.agents/skills/optimize-workflow/trials/[workflow]/[trial-2]/report.json --out-file .github/.agents/skills/optimize-workflow/trials/[workflow]/[trial-1]-[trial-2]-comparison.md`
7. Update trial log. Include updated before-after diagram if structural changes occurred.
8. Show condensed results. Prompt for more trials or stop.

<trial-template>
| Runner | Structure | Experiment | Expectation | Branch | Run ID | Commit | Stability | Runtime | Cost | Notes |
|---|---|---|---|---|---|---|---|---|---|---|
| `config` | `default / cached / parallel` | What's tested | Expected result | `branch` | `run-id` | `sha` | Pass/Fail/Flaky | mm:ss | $ | Findings |
</trial-template>
</loop>

<complete trigger="User stops loop or all trials done">
1. Recommend final configuration.
2. Add PR summary table and before-after Mermaid diagrams:

```md
### [Workflow/Job Name] Runner & Structure Changes

| Approach | Runner | Structure | Stability | Runtime | Runtime Delta (Abs/%) | Cost | Cost Delta (Abs/%) |
|---|---|---|---|---|---|---|---|
| [Old](https://github.com/link/to/baseline/workflow_run) | [original runner] | [original structure] | Pass/Fail/Flaky | mm:ss | +0:00 (+0%) | $ | +$ (+0%) |
| [New](https://github.com/link/to/final/workflow_run) | [new runner] | [new structure] | Pass/Fail/Flaky | mm:ss | +0:00 (+0%) | $ | +$ (+0%) |
```

3. Clean debug lines. Make final approved edits on a new branch. Ask user to commit.
4. Delete trial branches, logs, and PRs.
</complete>