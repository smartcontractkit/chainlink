---
name: optimize-workflow
description: Optimize CI runners and workflow structure for speed, simplicity, stability, and cost.
disable-model-invocation: true
---

<initialization>
1. Confirm missing inputs: target workflow, scope, priorities (default: stability > speed > cost).
2. Setup baseline:
   - Analyze target workflow (runners, caching, bottlenecks, critical path).
   - Setup test workflow (add `workflow_dispatch`, mock inputs if needed).
   - Init/resume trial log: `.github/.agents/skills/optimize-workflow/trials/[workflow]/summary.md`.
   - Record baseline run metrics.
</initialization>

<git_workflow>
To prevent index/branch lock conflicts:
- Do NOT checkout local trial branches.
- Push trial commits directly from working branch to remote refs:
  `git push -f origin HEAD:refs/heads/trial/[trial-name]`
- Delete remote trial refs on cleanup: `git push origin --delete trial/[trial-name]`
</git_workflow>

<constraints>
- Git: Direct ref pushes only. No local branch switching during trials.
- Runners: Prefer `ubuntu-latest`. Use `runs-on` for larger runners (verify API first).
- Scope: One change per trial. Keep cache status identical across comparison runs.
- Validation: Lint YAML (`actionlint` or dry-run) before pushing.
</constraints>

<tools>
If not already installed, prompt user to install [octometrics](https://github.com/kalverra/octometrics)

```sh
octometrics -h # Help menu
# Get data on specific workflow run, PR, or commit
octometrics [url] --format [json|md] -f .github/.agents/skills/optimize-workflow/trials/[workflow]/[description-of-run].[json|md] --download-logs

# See cleaned up logs of job run
octometrics log [job_run_id]

# Compare two runs
octometrics [before-url] --vs [after-url] --format [json|md] -f .github/.agents/skills/optimize-workflow/trials/[workflow]/[description-of-comparison].[json|md]
```

Use `octometrics` first to gather data. If data is missing or incomplete, then use `gh` CLI.
</tools>

<resources>
- [runs-on docs](https://runs-on.com/docs/) | [spot pricing](https://runs-on.com/docs/costs/spot-pricing/)
- [available runners](https://go.runs-on.com/api)
- [GitHub Actions docs](https://docs.github.com/en/actions)
</resources>

<loop>
1. Define single-variable trial and prompt for user approval.
2. Push commit to remote ref: `git push -f origin HEAD:refs/heads/trial/[trial-name]`.
3. Trigger run: `gh workflow run --ref trial/[trial-name]`.
4. Monitor run with `workflow_monitor.py`.
5. Compare with `workflow_compare.py` against baseline/previous trial.
6. Update trial log in `summary.md`. Include before/after diagram if structure changed.
7. Present condensed result to user; prompt for next trial or completion.

<trial-template>
| Runner | Structure | Experiment | Expectation | Branch | Run ID | Commit | Stability | Runtime | Cost | Notes |
|---|---|---|---|---|---|---|---|---|---|---|
| `config` | `cached/parallel` | Tested change | Expected impact | `trial/...` | `id` | `sha` | Pass/Fail | mm:ss | $ | Notes |
</trial-template>
</loop>

<complete trigger="User stops loop or all trials done">
1. Recommend final configuration.
2. Output PR summary table and before/after Mermaid diagrams:

```md
### [Workflow/Job Name] Runner & Structure Changes

| Approach | Runner | Structure | Stability | Runtime | Runtime Delta (Abs/%) | Cost | Cost Delta (Abs/%) |
|---|---|---|---|---|---|---|---|
| [Old](link) | [runner] | [structure] | Pass/Fail | mm:ss | +0:00 (+0%) | $ | +$ (+0%) |
| [New](link) | [runner] | [structure] | Pass/Fail | mm:ss | +0:00 (+0%) | $ | +$ (+0%) |
```

3. Clean up debug code. Make final edits on working branch. Ask user to commit.
4. Delete remote trial branches (`git push origin --delete trial/[trial-name]`) and temporary files.
</complete>
