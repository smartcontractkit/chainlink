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
| `runs-on/cpu=16/ram=64/family=m7i+m8i/extras=s3-cache+tmpfs` | Trial 1 | faster non-RMN tests | — | — | — | — | — | — | pending |
| `runs-on/cpu=48/ram=192/family=m6i+m7i/extras=s3-cache+tmpfs` | Trial 2 (RMN only) | faster RMN if CPU/RAM bound | — | — | — | — | — | — | pending |
| `ubuntu-latest-16cores-64GB` | Trial 3 (build fallback) | fast GH fallback | — | — | — | — | — | — | pending |
