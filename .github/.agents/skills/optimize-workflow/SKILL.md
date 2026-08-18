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

<tools>
<octometrics>
Use [octometrics](https://github.com/kalverra/octometrics) CLI for workflow analysis and comparison.

```sh
octometrics -h # Help message
octometrics [url] --download-logs --format [json|md] -f [output-file].[json|md] # Gather workflow run/job run/commit/pr CI run metrics and logs
octometrics logs [job-run-id|job.run.url] # See job logs
octometrics [before-url] --vs [after-url] --format [json|md] -f [output-file].[json|md] # Compare workflow run/job run metrics and logs
```
</octometrics>
<available-runners>
Use [available-runners API](https://go.runs-on.com/api) to find valid `runs-on` runner configs.

```sh
curl -sG "https://go.runs-on.com/api/finder" [params] | jq [jq filters]
```
<params>
region - AWS region (e.g., us-east-1) or all (default: all, searches all regions). Multiple regions can be specified separated by comma or plus sign (e.g., us-east-1,us-west-2 or us-east-1+us-west-2). When multiple regions are specified, only instances available in all specified regions are returned.
cpu - CPU range in format min-max or single value (e.g., 2-4 or 4) or all (default: all)
ram - RAM in GB, format min-max or single value (e.g., 8-16 or 8) or all (default: all)
gpu - GPU count range in format min-max or single value (e.g., 1-2 or 1) or all (default: all)
passmark - PassMark score range (optional, format: min-max or single value)
arch - Architecture filter (optional, can be multiple): arm64, amd64, or x86_64
platform - Platform filter (optional, can be multiple): linux, linux/unix, or windows
family - Instance family or partial instance type filter (optional, comma or plus separated): Supports exact family names (e.g., m5), partial matching (e.g., m7 matches m7a, m7g, etc.), wildcard patterns (e.g., m5.* or m5.), or exact instance types (e.g., m5.large). Multiple families can be specified with commas or plus signs (e.g., m5,m7a or m5+m7a)
</params>
</available-runners>
</tools>

<constraints>
- Git: Use worktrees for local trial branches. Avoid local checkout of trial branches to prevent index/branch lock conflicts.
- Runners: Prefer `ubuntu-latest`. Use `runs-on` for larger runners. Verify all `runs-on` runners through <available-runners> before use.
- Scope: One change per trial. Keep cache status identical across comparison runs.
- Validation: Lint YAML (`actionlint` or dry-run) before pushing.
- Tools: Use `octometrics` for analysis and data gathering. Only use `gh` CLI as a backup when `octometrics` is insufficient
</constraints>

<resources>
- [runs-on docs](https://runs-on.com/docs/)
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
