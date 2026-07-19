# Trial Log: integration-tests.yml CCIP Shards Runner Optimization

Optimize runner sizes for CCIP test shards from `cpu=16/ram=64` to cheaper alternatives.

## Trial History

| Runner | Structure | Experiment | Expectation | Branch | Run ID | Commit | Stability | Runtime | Cost | Notes |
|---|---|---|---|---|---|---|---|---|---|---|
| `cpu=16/ram=64/family=m7i+m8i` | `parallel` | Baseline from previous runs | Run tests successfully, reference performance | `optimizeIntegrationTests` | `29658759993` | `6c1de1b` | Pass | 17:26 | $0.27/hr | Memory usage < 10GB. CPU load average average ~1.0. |
| `cpu=8/ram=32/family=m7i+m8i` | `parallel` | Reduce resources to 8 vCPUs and 32GB RAM | Run successfully, no resource failures | `trial/ccip-cpu8-ram32` | `29689709024` | `2e61f7c` | Fail | N/A | N/A | Terminated with exit code 143 (OOM) due to multiple Docker containers. |
| `cpu=8/ram=64/family=r7i` | `parallel` | Reduce CPU to 8 vCPUs, keep 64GB RAM | Run successfully, keep tmpfs, prevent OOM | `trial/ccip-cpu8-ram32` | `29690694050` | `927d64a` | Pass | 12:13 | $0.14/hr | Succeeded. Ran faster than baseline. ~50% cost savings. |
| `cpu=8/ram=32/family=m7i+m8i` | `parallel` | Remove tmpfs, keep 32GB RAM | Run successfully, save RAM by disabling tmpfs | `trial/ccip-cpu8-ram32` | `29696233411` | `aa8b06e` | Fail | N/A | N/A | Fails with connection timeouts/errors. Saturated EBS IOPS (disk bottlenecks) without tmpfs, dropping RPC connections. |

## Major Workflow Optimizations

1. **CCIP Runner Sizing**: Reduced CCIP E2E test shards from `cpu=16/ram=64` (m7i.4xlarge, ~$0.27/hr) to `cpu=8/ram=64/family=r7i` (r7i.2xlarge, ~$0.14/hr) with `tmpfs`. Attempts to use smaller memory (`ram=32`) with `tmpfs` fail due to OOM from docker node requirements. Attempts to run without `tmpfs` fail because EBS IOPS limits are exceeded, dropping Geth/Chainlink node RPC connections. Thus, `cpu=8/ram=64/family=r7i/extras=s3-cache+tmpfs` is the optimal runner configuration.
2. **Compile-Tests Fix**: Added `tmpfs` back to `compile-tests` job. Compilation on 32 cores without `tmpfs` saturated EBS IOPS, causing runner connection loss and timeout at 13m. Adding `tmpfs` resolved the disk bottleneck, reducing runtime to **1m 46s**.
3. **Decoupled Compile-Tests**: Conditionalized `compile-tests` to only run if CCIP E2E tests are scheduled (`should-run == 'true'`). Removed `compile-tests` dependency from Core CRE jobs, skipping unnecessary CCIP compilation when only Core tests run.

