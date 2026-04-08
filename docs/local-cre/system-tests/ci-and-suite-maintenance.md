---
id: local-cre-system-tests-ci-and-suite-maintenance
title: CI and Suite Maintenance
sidebar_label: CI and Suite Maintenance
sidebar_position: 3
---

# CI and Suite Maintenance

The smoke suite is designed to be discoverable and repeatable across multiple topologies.

## Auto-Discovery

At a high level, the CI flow works like this:

- CI discovers tests in `system-tests/tests/smoke/cre`
- it builds a matrix across supported topologies and test entrypoints
- each test runs with configuration appropriate to that topology

In practice, this means you do not need to manually register every new CRE smoke test in a bespoke list.

## Naming and Placement

For discoverability and consistency:

- put the test in `system-tests/tests/smoke/cre`
- use the `Test_CRE_` prefix
- follow the existing package patterns for setup and helper usage

## Architecture Pattern

The suite uses a separated pattern:

- environment creation happens once per topology
- multiple tests reuse that environment
- deployed contracts and nodes are shared within the topology run

This keeps the suite cheaper and faster than recreating Local CRE for every individual test.

## Parallelism and Shared Resources

Parallel execution is controlled by `CRE_TEST_PARALLEL_ENABLED`.

That flag only permits parallelism. Each test still decides whether it is safe to call `t.Parallel()`. In `cre_suite_test.go`:

- some scenarios parallelize immediately
- some scenarios parallelize only when `CRE_TEST_CHIP_SINK_FANOUT_ENABLED=1`
- some scenarios stay serial because they depend on non-shareable infrastructure

The fanout flag matters for tests that share the ChIP test sink. In fanout mode, the helper starts one sink server and distributes events to per-test subscribers. That keeps parallel cases isolated without forcing a separate sink process per test.

## Supported Topologies in CI

By default, the CRE workflow runs tests against:

- `workflow-gateway-capabilities`

Some tests must replace that default topology set with explicit per-test overrides in `.github/workflows/cre-system-tests.yaml`. Current examples are:

- `Test_CRE_V2_Aptos_Suite` -> `workflow-gateway-aptos`
- `Test_CRE_V2_Solana_Suite` -> `workflow`
- `Test_CRE_V1_Tron` -> `workflow`
- `Test_CRE_V2_Sharding` -> `workflow-gateway-sharded`

If a new test only works with a non-default topology, adding the test code is not enough. You must also add an explicit override in the workflow matrix so CI runs the test with the matching `topology` and `configs` pair.

Keep tests topology-agnostic where possible. Use a per-test topology override only when the workflow genuinely depends on a different chain family or topology layout.

## Adding a New Test

When adding a new smoke test:

1. place it in the smoke package
2. follow the existing naming convention
3. prefer shared helpers
4. decide whether it belongs in an existing bucket or needs a new entrypoint
5. if it needs a non-default topology, add the explicit workflow-matrix override
6. verify it works with the expected topology matrix
7. keep external dependencies explicit

## Troubleshooting Discovery

If CI does not pick up a test:

1. confirm the function name starts with `Test_`
2. confirm the file is in the smoke package
3. confirm the package compiles
4. confirm the test does not depend on local-only assumptions absent in CI
