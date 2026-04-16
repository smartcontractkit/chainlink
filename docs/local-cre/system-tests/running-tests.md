---
id: local-cre-system-tests-running-tests
title: Running CRE System Tests
sidebar_label: Running Tests
sidebar_position: 1
---

# Running CRE System Tests

This page covers local, Kubernetes, and CI-facing execution concerns for the CRE smoke suite.

## Local Run Flow

The usual local path is:

```bash
cd core/scripts/cre/environment
go run . env setup
go run . env start
```

Then run the tests:

```bash
go test ./system-tests/tests/smoke/cre -timeout 20m -run '^Test_CRE_'
```

The comments in `cre_suite_test.go` also call out the pattern of starting Local CRE first and then running the smoke package.

Do not enable `--with-beholder` for the default smoke-test flow. Most CRE smoke tests start the ChIP test sink on the default gRPC port (`50051`), and Beholder starts Chip Ingress on that same port. If both try to use the default port, the test sink startup fails.

Enable Beholder only when:

- you are running Beholder-specific coverage
- you need the Beholder stack for debugging
- you start it on a different `--grpc-port` so the smoke-test sink can still bind its default port

## Environment Variables

The main variables used by the smoke suite are:

- `CTF_CONFIGS` points to the topology TOML before startup and to the generated Local CRE state file after startup. The helpers in `system-tests/tests/test-helpers/before_suite.go` switch to the state file automatically for local runs.
- `TOPOLOGY_NAME` is used in test names, bucket labels, and log output so results stay tied to the topology under test.
- `CTF_LOG_LEVEL=debug` enables more verbose framework logs during setup and test execution.
- `CTF_JD_IMAGE` pins the Job Distributor image when you do not want the default local image selection.
- `CTF_CHAINLINK_IMAGE` pins the Chainlink node image. `system-tests/lib/cre/environment/dons.go` checks this variable directly when selecting node images.

## Parallel Execution

Parallel test execution is opt-in. Set:

```bash
CRE_TEST_PARALLEL_ENABLED=1
```

The smoke suite does not blindly parallelize every case. The runner in `cre_suite_test.go` enables `t.Parallel()` only for scenarios that are safe to run together.

Current behavior:

- tests such as Proof of Reserve, Vault DON, and HTTP Trigger Action can run in parallel when `CRE_TEST_PARALLEL_ENABLED=1`
- tests such as HTTP Action CRUD, DON Time, Consensus, EVM Write, and EVM LogTrigger require both parallel mode and ChIP sink fanout mode
- Cron Beholder stays serial because it uses the real ChIP ingress stack instead of the shared fanout helper

## ChIP Sink Fanout Mode

Some workflows depend on the ChIP test sink. Running those cases in parallel requires the fanout server mode:

```bash
CRE_TEST_PARALLEL_ENABLED=1 \
CRE_TEST_CHIP_SINK_FANOUT_ENABLED=1 \
go test ./system-tests/tests/smoke/cre -timeout 20m -run '^Test_CRE_V2'
```

`CRE_TEST_CHIP_SINK_FANOUT_ENABLED` starts a singleton sink and fans events out to per-test subscribers. That allows multiple tests to share the sink without one test consuming another test's events or closing shared channels underneath another run.

Use fanout mode when parallelizing:

- HTTP Action CRUD
- DON Time
- Consensus
- EVM Write
- EVM LogTrigger

Without fanout mode, those cases still run, but they stay serial even if `CRE_TEST_PARALLEL_ENABLED=1`.

## Topology Defaults

For the default local flow, use:

`core/scripts/cre/environment/configs/workflow-gateway-capabilities-don.toml`

Override `CTF_CONFIGS` only when you intentionally need a different topology such as sharded, gateway-only, or chain-specific coverage.

## Timeouts and Debugging

Practical timeout guidance:

- allow about 20 minutes when the image is built from source
- expect much shorter runs when using prebuilt images

For local debugging, a useful pattern is:

```bash
CTF_LOG_LEVEL=debug \
go test ./system-tests/tests/smoke/cre -timeout 20m -run '^Test_CRE_V2_Suite_Bucket_A$'
```

That keeps the run narrow while preserving the topology and workflow setup used by the full suite.

Example VS Code launch configuration:

```json
{
  "name": "Launch CRE V2 Bucket A",
  "type": "go",
  "request": "launch",
  "mode": "test",
  "program": "${workspaceFolder}/system-tests/tests/smoke/cre",
  "args": ["-test.run", "^Test_CRE_V2_Suite_Bucket_A$"]
}
```

## Bucketed Test Selection

The larger CRE smoke suites are split into runtime-balanced buckets instead of one oversized test entrypoint.

The old V2 suite is split into:

- `Test_CRE_V2_Suite_Bucket_A`
- `Test_CRE_V2_Suite_Bucket_B`
- `Test_CRE_V2_Suite_Bucket_C`

Those buckets are defined in `system-tests/tests/smoke/cre/v2suite/config/bucketing.go`:

- `suite-bucket-a`: ProofOfReserve, HTTPTriggerAction, DONTime, Consensus
- `suite-bucket-b`: VaultDON
- `suite-bucket-c`: CronBeholder, HTTPActionCRUD

The EVM read suite uses a separate bucket registry in `system-tests/tests/smoke/cre/evm/evmread/config/bucketing.go`:

- `Test_CRE_V2_EVM_Read_HeavyCalls`
- `Test_CRE_V2_EVM_Read_StateQueries`
- `Test_CRE_V2_EVM_Read_TxArtifacts`

Use bucketed entrypoints when you want:

- shorter local feedback loops
- more stable CI runtimes
- a controlled way to rebalance scenario runtimes as the suite grows

## Minimal Troubleshooting Checklist

If a test run fails before test logic starts:

1. confirm `CTF_CONFIGS`
2. confirm the Local CRE state file is valid
3. confirm required images exist
4. rerun `go run . env setup`
5. rerun with debug logging
