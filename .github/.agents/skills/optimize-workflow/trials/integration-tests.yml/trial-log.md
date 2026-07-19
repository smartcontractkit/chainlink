# Trial Log: integration-tests.yml CCIP Shards Runner Optimization

Optimize runner sizes for CCIP test shards from `cpu=16/ram=64` to cheaper alternatives.

## Trial History

| Runner | Structure | Experiment | Expectation | Branch | Run ID | Commit | Stability | Runtime | Cost | Notes |
|---|---|---|---|---|---|---|---|---|---|---|
| `cpu=16/ram=64/family=m7i+m8i` | `parallel` | Baseline from previous runs | Run tests successfully, reference performance | `optimizeIntegrationTests` | `29658759993` | `6c1de1b` | Pass | 17:26 | $ | Memory usage < 10GB. CPU load average average ~1.0. |
| `cpu=8/ram=32/family=m7i+m8i` | `parallel` | Reduce resources to 8 vCPUs and 32GB RAM | Run successfully, no resource failures, cheaper cost | `trial/ccip-cpu8-ram32` | TBD | TBD | TBD | TBD | TBD | Target trial |

## Proposed Trials

1. **Trial 1**: Change runner configuration for all CCIP E2E test shards in `.github/e2e-tests.yml` to use `cpu=8/ram=32` instead of `cpu=16/ram=64`.
2. **Trial 2**: If Trial 1 succeeds and shows no resource constraint failures, test `cpu=8/ram=16` (to run even cheaper).
