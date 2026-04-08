---
id: local-cre-system-tests-index
title: CRE System Tests
sidebar_label: System Tests
sidebar_position: 0
---

# CRE System Tests

The CRE smoke suite lives in `system-tests/tests/smoke/cre`. These tests are designed to run against a Local CRE environment and reuse the same workflow, topology, and registry concepts documented in the environment pages.

## Smoke vs Regression

The package-level rule is:

- smoke tests cover happy-path and sanity-check behavior
- edge cases and negative conditions belong in `system-tests/tests/regression/cre`

## How the Tests Use Local CRE

The helper flow is implementation-backed:

- if `CTF_CONFIGS` is empty, the helpers set it to the requested config
- if the Local CRE state file does not exist, the helpers start Local CRE with `go run . env start`
- after environment creation, `CTF_CONFIGS` is switched to the local CRE state file so the tests use the deployed environment rather than the original topology TOML

This is why you can either:

- start Local CRE manually, then run tests
- or let the helpers bootstrap it for you

## Main Topics

- [Running Tests](running-tests.md)
- [Workflows in Tests](workflows-in-tests.md)
- [CI and Suite Maintenance](ci-and-suite-maintenance.md)

## Recommended Local Flow

1. bring up Local CRE
2. confirm the desired topology
3. run the relevant smoke tests
4. use Beholder or observability only when the scenario needs extra signals
