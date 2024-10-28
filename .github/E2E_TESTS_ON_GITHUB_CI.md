# E2E Tests on GitHub CI

- [E2E Tests on GitHub CI](#e2e-tests-on-github-ci)
  - [Scheduled test workflows](#scheduled-test-workflows)
    - [PR E2E Tests](#pr-e2e-tests)
    - [Nightly E2E Tests](#nightly-e2e-tests)
    - [Release E2E Tests](#release-e2e-tests)
      - [Integration (smoke) Tests](#integration-smoke-tests)
      - [Client Compatibility Tests](#client-compatibility-tests)
  - [On-Demand Workflows](#on-demand-workflows)
    - [Test workflows setup in CI](#test-workflows-setup-in-ci)
      - [Configuration Overrides](#configuration-overrides)
      - [Test Secrets](#test-secrets)

E2E tests are executed on GitHub CI using the [E2E Tests Reusable Workflow](https://github.com/smartcontractkit/.github/blob/main/.github/workflows/README.md) or dedicated workflows.

## Scheduled test workflows

These workflows are designed to run on every commit in a PR, nightly or before release (see `triggers` in [e2e-tests.yaml](./e2e-tests.yml)).

### PR E2E Tests

Tests triggered on every commit in a PR to ensure changes do not introduce regressions.

**Workflow:** [integration-tests.yml](https://github.com/smartcontractkit/chainlink/blob/develop/.github/workflows/integration-tests.yml)

### Nightly E2E Tests

Nightly E2E test runs.

**Workflow:** [nightly-e2e-tests.yml](https://github.com/smartcontractkit/chainlink/blob/develop/.github/workflows/run-nightly-e2e-tests.yml)

### Release E2E Tests

E2E tests triggered on a release tag.

#### Integration (smoke) Tests

**Workflow:** [integration-tests.yml](https://github.com/smartcontractkit/chainlink/blob/develop/.github/workflows/integration-tests.yml)

#### Client Compatibility Tests

**Workflow:** [client-compatibility-tests.yml](https://github.com/smartcontractkit/chainlink/actions/workflows/client-compatibility-tests.yml)

## On-Demand Workflows

These are dispatched parametrized workflows, that may be triggered manually for specific testing needs. For more details refer [integration-tests README](../integration-tests/README.md) and [per-product test run books](../integration-tests/run-books/).

**Examples:**

- [Selected E2E Tests Workflow](https://github.com/smartcontractkit/chainlink/actions/workflows/run-selected-e2e-tests.yml)
- [Client Compatibility Tests](https://github.com/smartcontractkit/chainlink/actions/workflows/client-compatibility-tests.yml)
- [Chaos Tests](https://github.com/smartcontractkit/chainlink/actions/workflows/integration-chaos-tests.yml)
- [OCR Soak Tests](https://github.com/smartcontractkit/chainlink/actions/workflows/on-demand-ocr-soak-test.yml)
- [On-Demand Automation Tests](https://github.com/smartcontractkit/chainlink/actions/workflows/automation-ondemand-tests.yml)
- [CCIP Chaos Tests](https://github.com/smartcontractkit/chainlink/actions/workflows/ccip-chaos-tests.yml)
- [CCIP Load Tests](https://github.com/smartcontractkit/chainlink/actions/workflows/ccip-load-tests.yml)
- [VRFv2Plus Smoke Tests](https://github.com/smartcontractkit/chainlink/actions/workflows/on-demand-vrfv2plus-smoke-tests.yml)
- [VRFv2Plus Performance Tests](https://github.com/smartcontractkit/chainlink/actions/workflows/on-demand-vrfv2plus-performance-test.yml)

### Test workflows setup in CI

Most workflows may be triggered with default configs. Some, nevertheless, may be overridden.

> [!TIP]
> Use `gh` CLI commands to run workflows from local machine.

#### Configuration Overrides

> [!CAUTION]
> Test configurations should not keep any [sensitive data or secrets](#test-secrets).

1. Reference sources:
   1. [Integration-Tests configurations](../integration-tests/testconfig/README.md);
   2. [CTF Test Config](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/lib/config/README.md).
2. Defaults and overrides should be stored (committed) in repository under `../integration-tests/testconfig/<product>/overrides/<override-name>.toml` (see example [here](../integration-tests/testconfig/ocr2/overrides/base_sepolia.toml)).
3. Use `test_config_override_path` to point to an override config. For example: `test_config_override_path="testconfig/ocr2/overrides/base_sepolia.toml"`

#### Test Secrets

> [!CAUTION]
> Pay attention to never store/expose/commit your test secrets in repository.

Test secrets allow provisioning and override the sensitive data such as EOA's private key, RPCs, Docker registry links, etc.

Reference sources:

1. [BASE64_CONFIG_OVERRIDE](../integration-tests/testconfig/README.md#base64_config_override).
2. [CTF Test Secrets documentation](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/lib/config/README.md#test-secrets).
3. [Guide on running GitHub workflows with your test secrets](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/lib/config/README.md#run-github-workflow-with-your-test-secrets).
