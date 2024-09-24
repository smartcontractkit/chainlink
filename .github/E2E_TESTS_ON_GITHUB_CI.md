# E2E Tests on GitHub CI

E2E tests are executed on GitHub CI using the [E2E Tests Reusable Workflow](#about-the-reusable-workflow) or dedicated workflows.

## Automatic workflows

These workflows are designed to run automatically at crucial stages of the software development process, such as on every commit in a PR, nightly or before release.

### PR E2E Tests

Run on every commit in a PR to ensure changes do not introduce regressions.

[Link](https://github.com/smartcontractkit/chainlink/blob/develop/.github/workflows/integration-tests.yml)

### Nightly E2E Tests

Conducted nightly to catch issues that may develop over time or with accumulated changes.

[Link](https://github.com/smartcontractkit/chainlink/blob/develop/.github/workflows/run-nightly-e2e-tests.yml)

### Release E2E Tests

This section contains automatic workflows triggering E2E tests at release.

#### Client Compatibility Tests

[Link](https://github.com/smartcontractkit/chainlink/actions/workflows/client-compatibility-tests.yml)

## On-Demand Workflows

Triggered manually by QA for specific testing needs.

**Examples:**

- [On-Demand Automation Tests](https://github.com/smartcontractkit/chainlink/actions/workflows/automation-ondemand-tests.yml)
- [CCIP Chaos Tests](https://github.com/smartcontractkit/chainlink/actions/workflows/ccip-chaos-tests.yml)
- [OCR Soak Tests](https://github.com/smartcontractkit/chainlink/actions/workflows/on-demand-ocr-soak-test.yml)
- [VRFv2Plus Smoke Tests](https://github.com/smartcontractkit/chainlink/actions/workflows/on-demand-vrfv2plus-smoke-tests.yml)

## Test Configs

E2E tests utilize TOML files to define their parameters. Each test is equipped with a default TOML config, which can be overridden by specifying an alternative TOML config. This allows for running tests with varied parameters, such as on a non-default blockchain network. For tests executed on GitHub CI, both the default configs and any override configs must reside within the git repository. The `test_config_override_path` workflow input is used to provide a path to an override config.

Config overrides should be stored in `testconfig/*/overrides/*.toml`. Placing files here will not trigger a rebuild of the test runner image.

**Important Note:** The use of `base64Config` input is deprecated in favor of `test_config_override_path`. For more details, refer to [the decision log](https://smartcontract-it.atlassian.net/wiki/spaces/TT/pages/927596563/Storing+All+Test+Configs+In+Git).

To learn more about test configs see [CTF Test Config](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/lib/config/README.md).

## Test Secrets

For security reasons, test secrets and sensitive information are not stored directly within the test config TOML files. Instead, these secrets are securely injected into tests using environment variables. For a detailed explanation on managing test secrets, refer to our [Test Secrets documentation](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/lib/config/README.md#test-secrets).

If you need to run a GitHub workflow using custom secrets, please refer to the [guide on running GitHub workflows with your test secrets](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/lib/config/README.md#run-github-workflow-with-your-test-secrets).

## About the E2E Test Reusable Workflow

For information on the E2E Test Reusable Workflow, visit the documentation in the [smartcontractkit/.github repository](https://github.com/smartcontractkit/.github/blob/main/.github/workflows/README.md).
