---
name: properly-size-runners
description: >-
  Run controlled Trials to size RunsOn self-hosted GitHub Actions runners for
  stability, speed, and cost. Use when tuning runs_on_self_hosted labels,
  diagnosing runner shutdown/OOM/spot kills, or comparing instance families
  (c/m/r, flex, tmpfs, spot) for e2e or CI jobs.
disable-model-invocation: true
---

# Properly Size Runners

Experiment-driven sizing for [RunsOn](https://runs-on.com/docs/) labels in this repo (typically `runs_on_self_hosted` in `.github/e2e-tests.yml` or workflow `runs-on:`).

Goal: pick `cpu` / `ram` / `family` / `spot` / `extras` that maximize **stability**, then trade **speed** vs **cost** per user priorities.

## Hard rules

1. Never change production runner labels without a recorded Trial result (or explicit user override).
2. Ask priorities before designing the experiment matrix.
3. Change **one primary lever** per Trial when possible (`ram`, `tmpfs`, `spot`, `family`, `cpu`).
4. Prefer evidence from job logs + metrics over intuition.
5. Cap concurrent expensive Trials; prefer `workflow_dispatch` / labeled PR runs over hammering develop.

## Inputs (ask if missing)

| Input | Required | Notes |
|-------|----------|-------|
| Target jobs / test IDs | yes | e.g. matrix entry in `.github/e2e-tests.yml` |
| Config path | yes | usually `.github/e2e-tests.yml` or workflow YAML |
| Baseline label | yes | current `runs_on_self_hosted` / `runs-on` |
| Priorities | yes | ranked: stability > speed > cost (or user order) |
| Budget | no | max $/run or max parallel Trials |
| Success criteria | no | defaults below |
| Tools | no | `gh`, octometrics, RunsOn metrics S3, CloudWatch |

### Default success criteria

- **Stability**: no runner shutdown / exit 143 / OOM during setup or test; job reaches test completion (PASS/FAIL/SKIP ok).
- **Speed**: wall-clock job duration (prefer p50 of ≥2 runs when flaky).
- **Cost**: instance type × duration × lifecycle (spot vs on-demand); prefer cheaper family when stability/speed within ~10%.

## Priority interview

Ask once at start. Record answers in the Trial log.

```text
1. Rank stability / speed / cost (1–3).
2. Acceptable failure modes? (spot retry ok? quarantined SKIP ok?)
3. Keep tmpfs? (yes / no / only if ram≥64)
4. Spot preference? (false | capacity-optimized | default)
5. Max Trials this session?
6. Apply winning label to config when done? (yes / no / PR only)
```

If user already stated priorities in chat, skip re-ask; confirm briefly.

## Failure mode cheat sheet

| Symptom | Likely cause | First lever |
|---------|--------------|-------------|
| Shutdown during gomod/cache `tar -xf` | tmpfs + low RAM | `ram=` ≥32 (match GH-hosted intent); drop `tmpfs` |
| Shutdown / exit 143 mid docker e2e | tmpfs docker + RAM pressure, or spot | drop `tmpfs`; `spot=false` or `capacity-optimized`; bump `ram` / `r*` |
| Shutdown early, identical labels on matrix | runner stealing | unique `runs-on=${{ github.run_id }}-<id>-${{ github.run_attempt }}` |
| Spot `BidEvictedEvent` | spot reclaim | `spot=false` or capacity-optimized + retry |
| Cancel + shutdown across jobs | workflow concurrency cancel | not a sizing issue |

EC2 naming (quick): `c`=compute, `m`=general, `r`=memory; `i`/`a`/`g`=Intel/AMD/Graviton; `flex`=cheaper bursty CPU. Always pin `ram=` when memory matters — `cpu=8` alone can land on 16 GiB `c*`.

Refs: [RunsOn labels](https://runs-on.com/docs/runners/labels/), [tmpfs](https://runs-on.com/docs/runners/capabilities/yolo-mode/), [Instance Finder](https://go.runs-on.com/finder), [rightsizing](https://runs-on.com/docs/performance/rightsizing/).

## Workflow

Copy and update:

```text
Sizing Progress:
- [ ] 0 Priorities confirmed
- [ ] 1 Baseline captured
- [ ] 2 Experiment matrix designed
- [ ] 3 Trials executed
- [ ] 4 Results compared
- [ ] 5 Recommendation approved
- [ ] 6 Config applied (if requested)
```

### 0. Priorities

Run priority interview. Derive scoring weights:

| Rank 1 | Rank 2 | Rank 3 | Weights (stability/speed/cost) |
|--------|--------|--------|--------------------------------|
| stability | speed | cost | 0.5 / 0.3 / 0.2 |
| stability | cost | speed | 0.5 / 0.2 / 0.3 |
| speed | stability | cost | 0.35 / 0.45 / 0.2 |
| cost | stability | speed | 0.35 / 0.2 / 0.45 |

Stability weight never below 0.35 unless user explicitly accepts flaky runners.

### 1. Baseline

Collect:

1. Current label string + resolved instance from last run (`InstanceType`, `InstanceRAM`, `InstanceLifecycle`, `Extras`).
2. Failure signature (step name, exit code, cache size if relevant).
3. Job URL(s) and run id.
4. Whether job is docker / in-memory / build; whether `tmpfs` / `s3-cache` / `spot` set.
5. GH-hosted fallback size if present (`ubuntu24.04-8cores-32GB` ⇒ aim `cpu=8/ram=32`).

Commands (prefix with `rtk` per workspace rules):

```bash
rtk gh run view <run_id> --repo <owner/repo> --json jobs
rtk gh run view <run_id> --repo <owner/repo> --job <job_id> --log
```

Extract from logs: `Runner details`, `InstanceType`, `InstanceRAM`, `##[error]`, `exit code`, `Cache Size`.

Optional: octometrics / `collect_test_telemetry: true` / RunsOn metrics S3 `metrics.jsonl` for memory/CPU peaks.

### 2. Design experiment matrix

Propose 2–5 Trials. One primary lever each.

**Starter ladders** (pick relevant):

| Trial | When | Label delta |
|-------|------|-------------|
| A baseline | always | current |
| B ram-fix | OOM / cache extract kill | add/raise `ram=32` (or 64), prefer `m*` over `c*` |
| C no-tmpfs | docker e2e + tmpfs + mid-test kill | remove `tmpfs` from `extras` |
| D spot-off | spot / unexplained kill | `spot=false` |
| E spot-co | cost after D stable | `spot=capacity-optimized` |
| F memory-class | docker + tmpfs kept | `family=r8i+r7i`, `ram=64`, keep or drop tmpfs per priority |
| G cpu-up | CPU-bound, stable RAM | `cpu=16` / `ram=64`, `m*` |

Use [Instance Finder](https://go.runs-on.com/finder) to sanity-check vCPU/RAM/region before proposing.

**Label template:**

```yaml
runs_on_self_hosted: runs-on/cpu=<N>/ram=<GiB>/family=<families>/spot=<policy>/image=ubuntu24-full-x64/extras=<extras>
```

Prefer unique runner prefix in workflows when matrix jobs share labels (fix in `run-e2e-tests` / calling workflow, not only static YAML).

Present matrix to user; wait for approval before mutating config or firing runs.

### 3. Execute Trials

For each approved Trial:

1. Patch only the Trial label (branch / local commit / dispatch input — follow user preference).
2. Trigger minimal run that exercises the target job(s).
3. Wait for completion (poll `gh run watch` / AwaitShell only if blocked).
4. Record results in Trial log (schema below).
5. Do not start next expensive Trial until prior result logged, unless user asked for parallel A/B.

If Trial fails infrastructure (shutdown/143) before test verdict → mark **unstable**, do not score speed/cost as wins.

### 4. Compare

Normalize per Trial:

| Field | Source |
|-------|--------|
| stable | bool — completed without runner kill |
| duration_s | job `completed_at - started_at` or step wall time |
| instance | `InstanceType` |
| ram_mib | `InstanceRAM` |
| lifecycle | spot / on-demand |
| extras | parsed from label |
| test_verdict | PASS / FAIL / SKIP / NONE |
| notes | failure signature |

Score (only if `stable`; else score = 0):

```text
score = w_s * 1.0
      + w_speed * (1 - duration_s / max_duration_among_stable)
      + w_cost  * (1 - est_cost / max_est_cost_among_stable)
```

`est_cost` ≈ duration_h × on-demand_or_spot_hourly for instance type (Finder / AWS price; rough OK).

Present a comparison table + winner under stated weights. Call out Pareto options (e.g. cheapest stable vs fastest stable).

### 5. Recommend

Output:

1. Winning label (copy-paste YAML).
2. Why (2–4 bullets tied to priorities + evidence).
3. Rejected alternatives and why.
4. Residual risks (flex CPU baseline, spot, matrix stealing, quarantine masking).
5. Ask: apply to config now?

### 6. Apply

If approved:

1. Update all relevant matrix entries (don't leave mixed labels unless intentional).
2. Keep GH-hosted `runs_on` fallback aligned in intent (`ram` ≈ cores×4 for m-class).
3. Summarize diff; commit only if user asks.

## Trial log schema

Append one block per Trial:

```markdown
### Trial <id> — <short name>
- label: `runs-on/...`
- primary_lever: <ram|tmpfs|spot|family|cpu|other>
- run_url: <url>
- job_id: <id>
- instance: <type> (<ram_mib> MiB, <lifecycle>)
- stable: <yes|no>
- duration_s: <n|n/a>
- test_verdict: <PASS|FAIL|SKIP|NONE>
- failure_signature: <none|shutdown@step|exit 143|oom|spot|other>
- metrics: <none|octometrics|cw|s3 path>
- notes: <free text>
```

## Experiment design anti-patterns

- Changing tmpfs + ram + spot in one Trial.
- Trusting quarantined SKIP as proof of docker-size correctness.
- Omitting `ram=` when comparing `c*` vs `m*`.
- Declaring cost winner from a killed run (short duration is not cheaper if unstable).
- Infinite re-runs without changing a lever.

## Local docs / tools

- `.github/AGENTS.md` — RunsOn preference, octometrics links
- `gh` for runs/jobs/logs
- [kalverra/octometrics](https://github.com/kalverra/octometrics) — workflow debugging
- [kalverra/octometrics-action](https://github.com/kalverra/octometrics-action) — resource monitoring

## Scaffold TODOs (fill in later)

- [ ] Script: parse RunsOn bootstrap table from job logs → JSON
- [ ] Script: score Trials from Trial log markdown / JSONL
- [ ] Default experiment presets per test_env_type (`docker` | `in-memory` | `build`)
- [ ] Wire `workflow_dispatch` inputs for label override without editing YAML
- [ ] Optional Cursor skill install path: copy/adapt to `.cursor/skills/properly-size-runners/SKILL.md`
