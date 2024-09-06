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

For security reasons, test secrets and sensitive information are not stored directly within the test config TOML files. Instead, these secrets are securely injected into tests using environment variables. For a detailed explanation on managing test secrets, refer to our [Test Secrets documentation](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/config/README.md#test-secrets).

If you need to run a GitHub workflow using custom secrets, please refer to the [guide on running GitHub workflows with your test secrets](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/config/README.md#run-github-workflow-with-your-test-secrets).

## About the Reusable Workflow

The [E2E Tests Reusable Workflow](https://github.com/smartcontractkit/chainlink/blob/develop/.github/workflows/run-e2e-tests-reusable-workflow.yml) is designed to run any type of E2E test on GitHub CI, including docker/testcontainers, old K8s tests, or tests in CRIB in the future.

Our goal is to migrate all workflows to use this reusable workflow for executing E2E tests. This approach will streamline our CI and allow for the automatic execution of tests at different stages of the software development process. Learn more about the advantages of using reusable workflows [here](https://smartcontract-it.atlassian.net/wiki/spaces/TT/pages/815497220/CI+Workflows+for+E2E+Tests).

**Examples of Workflows Utilizing the Reusable Workflow:**

- [Integration Tests](https://github.com/smartcontractkit/chainlink/blob/develop/.github/workflows/integration-tests.yml)
- [Nightly E2E Tests](https://github.com/smartcontractkit/chainlink/blob/develop/.github/workflows/run-nightly-e2e-tests.yml)
- [Selected E2E Tests](https://github.com/smartcontractkit/chainlink/blob/develop/.github/workflows/run-selected-e2e-tests.yml)
- [On-Demand Automation Tests](https://github.com/smartcontractkit/chainlink/blob/develop/.github/workflows/automation-ondemand-tests.yml)
- [CCIP Chaos Tests](https://github.com/smartcontractkit/chainlink/blob/develop/.github/workflows/ccip-chaos-tests.yml)

### E2E Test Configuration on GitHub CI

The [e2e-tests.yml](https://github.com/smartcontractkit/chainlink/blob/develop/.github/e2e-tests.yml) file lists all E2E tests configured for execution on CI. Each entry specifies the type of GitHub Runner needed and the workflows in which the test is included (for example, `smoke/ocr_test.go:*` is executed both in PRs and nightly).

### Slack Notification After Tests

To configure Slack notifications after tests executed via the reusable workflow, follow these steps:

- Set `slack_notification_after_tests` to either `always` or `on_failure` depending on when you want notifications to be sent.
- Assign `slack_notification_after_tests_channel_id` to the ID of the Slack channel where notifications should be sent.
- Provide a title for the notification by setting `slack_notification_after_tests_name`.
- Optionally use `slack_notification_after_tests_notify_user_id_on_failure` to reply in the thread and notify a user about the failed workflow

**Example:**

```yml
jobs:
  call-run-e2e-tests-workflow:
    name: Run E2E Tests
    uses: ./.github/workflows/run-e2e-tests-reusable-workflow.yml
    with:
      chainlink_version: develop
      test_workflow: Nightly E2E Tests
      slack_notification_after_tests: true
      slack_notification_after_tests_channel_id: "#team-test-tooling-internal"
      slack_notification_after_tests_name: Nightly E2E Tests
      slack_notification_after_tests_notify_user_id_on_failure: U0XXXXXXX
```

## Guides

### How to Run Selected Tests on GitHub CI

The [Run Selected E2E Tests Workflow](https://github.com/smartcontractkit/chainlink/actions/workflows/run-selected-e2e-tests.yml) allows you to execute specific E2E tests as defined in the [e2e-tests.yml](https://github.com/smartcontractkit/chainlink/blob/develop/.github/e2e-tests.yml). This is useful for various purposes such as testing specific features on a particular branch or verifying that modifications to a test have not introduced any issues.

For details on all available inputs that the workflow supports, refer to the [workflow definition](https://github.com/smartcontractkit/chainlink/actions/workflows/run-selected-e2e-tests.yml).

**Example Usage:**

To run a set of VRFv2Plus tests from a custom branch, use the following command:

```bash
gh workflow run run-selected-e2e-tests.yml \
-f test_ids="TestVRFv2Plus_LinkBilling,TestVRFv2Plus_NativeBilling,TestVRFv2Plus_DirectFunding,TestVRFV2PlusWithBHS" \
-f chainlink_version=develop \
--ref TT-1550-Provide-PoC-for-keeping-test-configs-in-git
```

### How to Run Custom Tests with Reusable Workflow

To run a specific list of tests, utilize the `custom_test_list_json` input. This allows you to provide a customized list of tests. If your test list is dynamic, you can generate it during a preceding job and then reference it using:  `custom_test_list_json: ${{ needs.gen_test_list.outputs.test_list }}`.  

```yml
run-e2e-tests-workflow:
  name: Run E2E Tests
  uses: ./.github/workflows/run-e2e-tests-reusable-workflow.yml
  with:
    custom_test_list_json: >
      {
        "tests": [
          {
            "id": "TestVRFv2Plus",
            "path": "integration-tests/smoke/vrfv2plus_test.go",
            "runs_on": "ubuntu-latest",
            "test_env_type": "docker",
            "test_cmd": "cd integration-tests/smoke && go test vrfv2plus_test.go"
          }
        ]
      }              
```