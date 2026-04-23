---
id: local-cre-index
title: Local CRE
sidebar_label: Overview
sidebar_position: 0
---

# Local CRE

Local CRE is the developer environment for building, running, and testing CRE workflows locally. It covers the environment lifecycle in `core/scripts/cre/environment` and the smoke-test flows in `system-tests/tests/smoke/cre`.

Use this doc set when you need to:

- bootstrap a local CRE stack on Docker or Kubernetes
- choose or inspect a topology
- deploy or debug workflows
- run or extend CRE smoke tests
- understand how the test helpers interact with a running Local CRE environment

## Start Here

- New to Local CRE: [Getting Started](getting-started/index.md)
- Running the environment day to day: [Environment](environment/index.md)
- Writing and debugging workflows: [Workflow Operations](environment/workflows.md)
- Choosing a topology: [Topologies and Capabilities](environment/topologies.md)
- Running or extending smoke tests: [System Tests](system-tests/index.md)
- Generated artifacts and references: [Reference](reference/index.md)

## Scope

This section is intentionally limited to Local CRE and the CRE system tests. It does not attempt to document the entire repository.

## Structure

This section is organized around the main Local CRE workflows:

- environment setup and lifecycle
- workflow deployment and debugging
- topology and capability selection
- CRE smoke-test execution and maintenance
