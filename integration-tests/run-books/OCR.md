
# OCR Tests Run-Book

- [OCR Tests Run-Book](#ocr-tests-run-book)
  - [Summary](#summary)
  - [Instructions](#instructions)
    - [Pre-requisites](#pre-requisites)
    - [Run Tests](#run-tests)
      - [COMMON COMMANDS](#common-commands)
        - [Docker](#docker)
        - [Kubernetes](#kubernetes)
        - [CI (with overrides)](#ci-with-overrides)
      - [SMOKE](#smoke)
        - [Docker](#docker-1)
        - [CI](#ci)
      - [SOAK](#soak)
        - [Kubernetes](#kubernetes-1)
        - [CI](#ci-1)
      - [LOAD](#load)
      - [CHAOS](#chaos)
        - [Kubernetes](#kubernetes-2)
        - [CI](#ci-2)
      - [MIGRATION (core version upgrade)](#migration-core-version-upgrade)
        - [Docker](#docker-2)

## Summary

This run-book is a guideline for running on demand OCR tests against any blockchain.

## Instructions

### Pre-requisites

> [!IMPORTANT]
>
> 1. Ensure [main Pre-requisites](../README.md#pre-requisites) are met.
> 2. Pay attention to the OCR version enabled in `overreride.toml`.
> 3. Use `-p 1` to disable tests parallelization and avoid nonce-related issues (or comment `t.Parallel()`).
> 4. For running tests in Kubernetes and CI, ensure test secrets are provided/uploaded to GitHub (ref. [CTF README#test-secrets](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/lib/config/README.md#test-secrets)).

### Run Tests

Below you may find instructions for running tests in different environments.

#### COMMON COMMANDS

Reuse the commands below to run different tests by their types/suites.

##### Docker

Any test suite/test can be run in Docker using the following `go` command (overrides are automatically injected):

```bash
go test -v -timeout <max_test_timeout> -p 1 <./path/to/_test.go file>
```

Example:

```bash
go test -v -timeout 60m -p 1 ./smoke/ocr2_test.go
```

##### Kubernetes

Run:

```bash
BASE64_CONFIG_OVERRIDE=$(cat ./testconfig/overrides.toml | base64) go test -v -timeout <max_test_timeout> -p 1 -run '<TestName>' ./<test_directory>
```

Example:

```bash
# Go
BASE64_CONFIG_OVERRIDE=$(cat ./testconfig/overrides.toml | base64) go test -v -p 1 -run 'TestOCRChaos' ./chaos

- - - - - - - - - - - - - -

# Make
BASE64_CONFIG_OVERRIDE=$(cat ./testconfig/overrides.toml | base64) make test_soak_ocr
```

##### CI (with overrides)

For the most tests run [Selected E2E Tests Workflow](https://github.com/smartcontractkit/chainlink/actions/workflows/run-selected-e2e-tests.yml) in GitHub, either manually or using `gh` command unless the otherwise stated:

```bash
gh workflow run "run-selected-e2e-tests.yml" \
--ref <branch> \
-f chainlink_version="v<your.version>" \ # Optional, default is image created from develop branch. Not needed if you run tests against existing environment
-f workflow_run_name="Any name" \ # Optional
-f test_ids="<test_ID>" \ # see /chainlink/.github/e2e-tests.yml for IDS
-f test_secrets_override_key=<your uploaded to GH secrets key> \ # Optional, can be obtained when secrets are uploaded to GitHub
-f test_config_override_path=<path/to/product/override.toml which should be committed in `integration-tests/testconfig/ocr2/overrides` folder> \ # Optional
-f with_existing_remote_runner_version=<remote-runner version> \ # Optional
```

Example:

```bash
gh workflow run "run-selected-e2e-tests.yml" \
--ref develop \
-f chainlink_version="v2.17.0-beta0" \
-f test_ids="smoke/ocr2_test.go:*" \
-f workflow_run_name="Smoke:OCR2:2.17.0-beta0" \
-f test_secrets_override_key=YOUR_TEST_SECRETS_ID \
-f test_config_override_path=./testconfig/ocr2/overrides/base_sepolia.toml
```

#### SMOKE

##### Docker

Refer [COMMON COMMANDS#Docker](#docker). Override `<./path/to/_test.go file>` as follows:

- With forwarders:
  `./smoke/forwarder_ocr_test.go`
  `./smoke/forwarder_ocr2_test.go`

- No forwarders:
  `./smoke/ocr_test.go`
  `./smoke/ocr2_test.go`

##### CI

Refer [COMMON COMMANDS#CI](#ci). Override `test_ids` as follows:

`smoke/ocr_test.go:*` - all OCR1 tests
`smoke/ocr2_test.go:*` - all OCR2 tests

#### SOAK

> [!IMPORTANT]
> These tests require logging in to Kubernetes cluster (`aws sso login`).
> Do not use `-timeout` flag for Soak tests with Go command. It is set in `overrides.toml`

Refer [Test config README](../testconfig/README.md) for more details about Soak tests configuration.

##### Kubernetes

Refer [COMMON COMMANDS#Kubernetes](#kubernetes). Override path as follows:

- With forwarders:
  `-run 'TestForwarderOCRv1Soak' ./soak` or `make test_soak_forwarder_ocr1`
  `-run 'TestForwarderOCRv2Soak' ./soak` or `make test_soak_forwarder_ocr2`

- No forwarders:
  `-run 'TestOCRv1Soak' ./soak` or `make test_soak_ocr`
  `-run 'TestOCRv2Soak' ./soak` or `make test_soak_ocr2`

- With reorg below finality and `FinalityTagEnabled=false`:
  `-run 'TestOCRSoak_GethReorgBelowFinality_FinalityTagDisabled' ./soak` or `make test_soak_ocr_reorg_1`

- With reorg below finality and `FinalityTagEnabled=true`:
  `-run 'TestOCRSoak_GethReorgBelowFinality_FinalityTagEnabled' ./soak` or `make test_soak_ocr_reorg_2`

- With gas spike:
  `-run 'TestOCRSoak_GasSpike' ./soak` or `make test_soak_ocr_gas_spike`

- With change of a block gas limit (creating block congestion):
  `-run 'TestOCRSoak_ChangeBlockGasLimit' ./soak` or `make test_soak_ocr_gas_limit_change`

- All RPCs get down for all nodes:
  `-run 'TestOCRSoak_RPCDownForAllCLNodes' ./soak` or `make test_soak_ocr_rpc_down_all_cl_nodes`

- 50% of nodes get RPCs down:
  `-run 'TestOCRSoak_RPCDownForHalfCLNodes' ./soak` or `make test_soak_ocr_rpc_down_half_cl_nodes`

##### CI

Use [On Demand OCR Soak Test](https://github.com/smartcontractkit/chainlink/actions/workflows/on-demand-ocr-soak-test.yml) workflow in GitHub.

OR

Run [On Demand OCR Soak Test](https://github.com/smartcontractkit/chainlink/actions/workflows/run-on-demand-ocr-soak-test.yml) workflow with `gh` as follows:

```bash
gh workflow run "on-demand-ocr-soak-test.yml" \
--ref <branch> \
-f chainlink_version="v<your.version>" \ # Optional, default is image created from develop branch. Not needed if you run tests against existing environment
-f testToRun="soak/ocr_test.go:<TestName>" \ # see /chainlink/.github/workflows/on-demand-ocr-soak-test.yml for options
-f test_secrets_override_key=<your uploaded to GH secrets key> \ # Optional, can be obtained when secrets are uploaded to GitHub
-f test_config_override_path=<path/to/product/override.toml which should be committed in `integration-tests/testconfig/ocr2/overrides` folder> \ # Optional
-f slackMemberID="YOUR_SLACK_MEMBER_ID" \ # Optional ("your profile -> three dots -> Copy Memeber ID")
```

Example:

```bash
gh workflow run "on-demand-ocr-soak-test.yml" \
--ref develop \
-f chainlink_version="v2.17.0-beta0" \
-f testToRun="soak/ocr_test.go:TestOCRv2Soak" \
-f test_config_override_path="/integration-tests/testconfig/ocr2/overrides/base_sepolia.toml" \
-f test_secrets_override_key=BASE_TESTSECRETS_YOUR_ID \ # RPC links in testsecret should correspond to the selected chain
-f slackMemberID="YOUR_SLACK_MEMBER_ID"
```

The following values may be used in the `testToRun` field (ref. [Soak#Kubernetes](#kubernetes) for more details):
`soak/ocr_test.go:TestOCRv1Soak`
`soak/ocr_test.go:TestOCRv2Soak`
`soak/ocr_test.go:TestForwarderOCRv1Soak`
`soak/ocr_test.go:TestForwarderOCRv2Soak`
`soak/ocr_test.go:TestOCRSoak_GethReorgBelowFinality_FinalityTagDisabled`
`soak/ocr_test.go:TestOCRSoak_GethReorgBelowFinality_FinalityTagEnabled`
`soak/ocr_test.go:TestOCRSoak_GasSpike`
`soak/ocr_test.go:TestOCRSoak_ChangeBlockGasLimit`
`soak/ocr_test.go:TestOCRSoak_RPCDownForAllCLNodes`
`soak/ocr_test.go:TestOCRSoak_RPCDownForHalfCLNodes`

#### LOAD

Ref: [Load Tests README](../load/ocr/README.md)

#### CHAOS

> [!IMPORTANT]
> 1. These tests require logging in to Kubernetes cluster (`aws sso login`).
> 2. There are only OCR1 chaos tests.

##### Kubernetes

Refer [COMMON COMMANDS#Kubernetes](#kubernetes). Override path as follows:

`-run 'TestOCRChaos' ./chaos`

OR

`make test_ocr_chaos`

##### CI

Refer [COMMON COMMANDS#CI](#ci-with-overrides). Override `test_ids` as follows:

`chaos/ocr_chaos_test.go`

#### MIGRATION (core version upgrade)

##### Docker

1. Refer [Test configurations README](../testconfig/README.md#migration-tests) to provide necessary configuration.
2. Run tests:

    ```bash
    # Go
    BASE64_CONFIG_OVERRIDE=$(cat ./testconfig/overrides.toml | base64) go test -v -p 1 ./smoke/<product>_test.go

    - - - - - - - - - - - - 

    # Make command
    BASE64_CONFIG_OVERRIDE=$(cat ./testconfig/overrides.toml | base64) make test_node_migrations_verbose
    ```
