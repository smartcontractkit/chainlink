# Right-Size Runners: Integration Tests (integration-tests.yml)

## Workflow

`integration-tests.yml` orchestrates the end-to-end test suite for the Chainlink monorepo.
The actual test definitions (runners, commands, timeouts) live in `e2e-tests.yml`, which is
consumed by the reusable `smartcontractkit/.github/.github/workflows/run-e2e-tests.yml`.

### Critical path in this workflow

```text
changes → labels
          ├─ run-core-cre-e2e-tests-setup ── build-chainlink ── run-core-cre-e2e-tests
          │                                                    └─ run-core-cre-e2e-regression-tests
          └─ run-ccip-v1-6-e2e-tests-setup ── build-chainlink ── run-ccip-v1-6-e2e-tests
                                                         └─ check-e2e-test-results
```

`build-chainlink` is a join point: it needs both the CRE and CCIP setup jobs to finish.
Its wall-clock is therefore gated by the slower of those two setup jobs. The actual CCIP
v1.6 tests are run as a matrix inside the reusable `run-e2e-tests.yml` workflow, so the
`run-ccip-v1-6-e2e-tests` job duration equals the slowest test in the filtered matrix.

## Current runner configuration

### `integration-tests.yml`

| Job | Runner | Notes |
|-----|--------|-------|
| `enforce-ctf-version` | `ubuntu-latest` | gate only |
| `changes` | `ubuntu-latest` | path filtering |
| `labels` | `ubuntu-latest` | resolves runner labels |
| `build-chainlink` (core) | `runs-on=${run_id}-core/cpu=32+36/memory=64+72/family=c6i+c7i+c5.*/extras=s3-cache+tmpfs` | heavy self-hosted |
| `build-chainlink` (plugins) | `runs-on=${run_id}-plugins/cpu=32+36/memory=64+72/family=c6i+c7i+c5.*/extras=s3-cache+tmpfs` | heavy self-hosted |
| `run-core-cre-e2e-tests-setup` | `ubuntu-latest` | gates only |
| `run-ccip-v1-6-e2e-tests-setup` | `ubuntu-latest` | gates only |
| `run-ccip-v1-6-e2e-tests` | caller (reusable workflow) | runner set inside `e2e-tests.yml` |
| `check-e2e-test-results` | `ubuntu-latest` | aggregation |

### `e2e-tests.yml` CCIP v1.6 entries

| Test ID | GH runner `runs_on` | Self-hosted `runs_on_self_hosted` | Timeout | Notes |
|---------|---------------------|-------------------------------------|---------|-------|
| `smoke/ccip/ccip_reorg_test.go:LessThanFinalityTests` | `ubuntu-latest` | — | 25m | no self-hosted runner defined |
| `smoke/ccip/ccip_reorg_test.go:GreaterThanFinalityTests` | `ubuntu-latest` | — | 25m | no self-hosted runner defined |
| `smoke/ccip/ccip_token_price_updates_test.go:*` | `ubuntu-latest` | — | 18m | `free_disk_space: true` |
| `smoke/ccip/ccip_gas_price_updates_test.go:^Test_CCIPGasPriceUpdatesWriteFrequency$` | `ubuntu-latest` | — | 18m | no self-hosted runner defined |
| `smoke/ccip/ccip_gas_price_updates_test.go:^Test_CCIPGasPriceUpdatesDeviation$` | `ubuntu-latest` | — | 18m | no self-hosted runner defined |
| `smoke/ccip/ccip_rmn_test.go:^TestRMN_TwoMessagesOneSourceChainCursed$` | `ubuntu24.04-16cores-64GB` | `runs-on/cpu=32/ram=128/family=m6i+m5.*/spot=false/image=ubuntu24-full-x64/extras=s3-cache+tmpfs` | 30m | disabled for Push/Nightly |
| `smoke/ccip/ccip_rmn_test.go:^TestRMN_GlobalCurseTwoMessagesOnTwoLanes$` | `ubuntu24.04-16cores-64GB` | `runs-on/cpu=32/ram=128/family=m6i+m5.*/spot=false/image=ubuntu24-full-x64/extras=s3-cache+tmpfs` | 30m | heaviest test |
| `smoke/ccip/ccip_jobspec_test.go.go:TestDeleteCCIPJobs` | `ubuntu-latest` | — | default | `free_disk_space: true` |
| `smoke/ccip/ccip_jobspec_test.go.go:TestRevokeJobs` | `ubuntu-latest` | — | default | `free_disk_space: true` |

The reusable runner-selection expression is:

```yaml
runs-on: ${{ inputs.use-self-hosted-runners == 'true' &&
  matrix.tests.runs_on_self_hosted || matrix.tests.runs_on }}
```

`integration-tests.yml` passes `use-self-hosted-runners: ${{ needs.labels.outputs.should-use-self-hosted-runners }}`,
which defaults to `true` unless the `runs-on-opt-out` label is present. Therefore, for the
majority of runs, every CCIP v1.6 test **without** a `runs_on_self_hosted` entry falls back
to `ubuntu-latest` (4 vCPU / 16 GiB), while the RMN tests get the large 32 vCPU / 128 GiB
self-hosted runner. That mismatch is the likely bottleneck: the non-RMN CCIP tests are
Docker-heavy but run on the smallest available runner.

## Optimization priorities (user-defined)

1. **Speed** > **Cost** > **Stability**.
2. "Speed" means reducing the total wall-clock time of `run-ccip-v1-6-e2e-tests`.
3. "Stability" allows some test flakes; focus is not resource-induced OOMs, which remain a hard constraint.
4. Scope: whole `integration-tests` workflow, with CCIP v1.6 as the stated bottleneck.

## Proposed trials

All trials are driven through `workflow_dispatch` on `integration-tests.yml` so they hit
the CCIP path without requiring label gymnastics or merge-group entry. The trial diff is
limited to `e2e-tests.yml` (adding/updating `runs_on_self_hosted`) unless the trial is
specifically targeting the build job.

| # | Runner config | Experiment | Expectation | Notes |
|---|---------------|------------|-------------|-------|
| 1 | Add `runs_on_self_hosted: runs-on/cpu=16/ram=64/family=m7i+m8i/extras=s3-cache+tmpfs` to all non-RMN CCIP v1.6 tests | Give CCIP smoke tests the same class of runner as the RMN GH fallback (16/64) so they stop running on `ubuntu-latest` | Non-RMN CCIP tests finish faster; wall-clock of `run-ccip-v1-6-e2e-tests` may drop if a non-RMN test was the critical path | Keeps `runs_on` unchanged; safe rollback by removing `runs_on_self_hosted` |
| 2 | Trial 1 + bump RMN `runs_on_self_hosted` to `cpu=48/ram=192/family=m6i+m7i/extras=s3-cache+tmpfs` | Test if the 30m RMN tests are CPU/memory bound | If RMN is the wall-clock, total runtime drops; if not, cost increases with no benefit | Only valuable if Trial 1 still leaves RMN as the slowest job |
| 3 | Trial 1 + set `build-chainlink` fallback to `ubuntu-latest-16cores-64GB` when `runs-on-opt-out` label is applied | Ensure the GH-hosted fallback path is also fast enough | Reduces risk for external contributors / opt-out runs; measures cost vs. self-hosted | Modifies `integration-tests.yml` labels job |

### Trial 1 detail

- `ccip_reorg_test.go:LessThanFinalityTests` → `runs_on_self_hosted: runs-on/cpu=16/ram=64/family=m7i+m8i/extras=s3-cache+tmpfs`
- `ccip_reorg_test.go:GreaterThanFinalityTests` → same
- `ccip_token_price_updates_test.go:*` → same
- `ccip_gas_price_updates_test.go:^Test_CCIPGasPriceUpdatesWriteFrequency$` → same
- `ccip_gas_price_updates_test.go:^Test_CCIPGasPriceUpdatesDeviation$` → same
- `ccip_jobspec_test.go.go:TestDeleteCCIPJobs` → same
- `ccip_jobspec_test.go.go:TestRevokeJobs` → same

This matches the RAM of the RMN GH fallback (`ubuntu24.04-16cores-64GB`) and provides a
large CPU uplift over `ubuntu-latest`. `m7i`/`m8i` are Intel general-purpose families
available in `us-east-1` per the runs-on finder API.

## Trial tracking table

| Runner | Experiment | Expectation | Workflow Run ID | Commit | Stability | Runtime | Cost | Score | Notes |
|--------|------------|-------------|-----------------|--------|-----------|---------|------|-------|-------|
| `ubuntu-latest` (baseline) | current non-RMN CCIP runner | reference | — | — | — | — | — | — | read-only baseline from recent runs |
| `runs-on/cpu=16/ram=64/family=m7i+m8i/extras=s3-cache+tmpfs` | Trial 1 | faster non-RMN tests | `29343491817` | `ea43a31c` | failed (RMN pre-existing timeout) | 32m46s matrix wall-clock; 13m37s non-RMN critical path | higher than `ubuntu-latest` | TBD pending baseline | see findings below |
| `runs-on/cpu=48/ram=192/family=m6i+m7i/extras=s3-cache+tmpfs` | Trial 2 (RMN only) | faster RMN if CPU/RAM bound | — | — | — | — | — | — | pending |
| `ubuntu-latest-16cores-64GB` | Trial 3 (build fallback) | fast GH fallback | — | — | — | — | — | — | pending |

## Trial 1 findings (run 29343491817)

### Overall result
- **Workflow run:** `failure`
- **Failure cause:** `smoke/ccip/ccip_rmn_test.go:^TestRMN_TwoMessagesOneSourceChainCursed$` hit its 30m test timeout (`panic: test timed out after 30m0s`).
- **Not caused by this trial:** that RMN test uses the unchanged `runs-on/cpu=32/ram=128/family=m6i+m5.*/spot=false/image=ubuntu24-full-x64/extras=s3-cache+tmpfs` runner documented as disabled for Push/Nightly. The other RMN test on the same runner passed in 6m49s.
- **ETH Smoke Tests** job failed only because it aggregates the CCIP result; core CRE tests passed.

### Runner verification
All non-RMN CCIP matrix jobs started on the intended runner:
```text
Labels            [runs-on=29343491817/cpu=16/ram=64/family=m7i+m8i/extras=s3-cache+tmpfs]
RunnerName        runs-on--i-02d0c2769adfbd027--i05u6zp859
ImageName         runs-on-v2.2-ubuntu24-full-x64-20260710155642
```

### Per-test durations (matrix jobs)

| Test | Runner | Duration | Result | Notes |
|------|--------|----------|--------|-------|
| `TestRMN_TwoMessagesOneSourceChainCursed` | 32 vCPU / 128 GiB (unchanged) | 32m46s | failure | 30m timeout; pre-existing issue |
| `ccip_reorg_test.go:GreaterThanFinalityTests` | 16 vCPU / 64 GiB (trial) | 13m37s | success | slowest non-RMN test |
| `ccip_reorg_test.go:LessThanFinalityTests` | 16 vCPU / 64 GiB (trial) | 10m55s | success | within 25m timeout |
| `Test_CCIPGasPriceUpdatesWriteFrequency` | 16 vCPU / 64 GiB (trial) | 7m49s | success | within 18m timeout |
| `TestRMN_GlobalCurseTwoMessagesOnTwoLanes` | 32 vCPU / 128 GiB (unchanged) | 6m49s | success | heaviest RMN test |
| `TestRevokeJobs` | 16 vCPU / 64 GiB (trial) | 4m51s | success | |
| `TestDeleteCCIPJobs` | 16 vCPU / 64 GiB (trial) | 4m34s | success | |
| `Test_CCIPGasPriceUpdatesDeviation` | 16 vCPU / 64 GiB (trial) | 2m57s | success | |
| `ccip_token_price_updates_test.go:*` | 16 vCPU / 64 GiB (trial) | 2m50s | success | within 18m timeout |

### Critical path
- `run-ccip-v1-6-e2e-tests` matrix wall-clock: **32m46s** (gated by the failing RMN test).
- Without the failing RMN test, the matrix wall-clock would have been **13m37s** (`GreaterThanFinalityTests`).
- The trial change did not affect the `build-chainlink` join point (6m setup time before CCIP jobs start).

### Stability assessment
- No OOM or resource-induced failures on the 16 vCPU / 64 GiB runner for non-RMN tests.
- All non-RMN CCIP tests passed on the new runner.
- The single failure is a pre-existing RMN timeout unrelated to the runner sizing change.

### Missing baseline
No recent `workflow_dispatch` run on `develop` executed the non-RMN CCIP tests on `ubuntu-latest`, so a direct speedup number is not available yet. PR/push runs always skip the CCIP matrix.

### Recommendations
1. **Run a baseline:** Create a temporary commit reverting the non-RMN `runs_on_self_hosted` entries so the same tests run on `ubuntu-latest`, then dispatch the workflow again. This gives a direct before/after comparison.
2. **Treat RMN timeout separately:** Investigate whether `TestRMN_TwoMessagesOneSourceChainCursed` is genuinely resource-bound or a test bug. If resource-bound, Trial 2 (bump RMN to 48/192) is the next step. If it is a flaky/disabled test, consider excluding it from workflow_dispatch or increasing its test timeout.
3. **Keep Trial 1 if baseline is favorable:** If the baseline shows non-RMN tests taking significantly longer than 13m37s on `ubuntu-latest`, the 16/64 runner is a clear win for the critical path.
